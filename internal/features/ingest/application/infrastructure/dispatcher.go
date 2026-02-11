package infrastructure

import (
	"sync"
	"time"

	"github.com/Jehoi-ga-ada/axiom-ingest-gateway/internal/features/ingest/domain"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

type DispatcherConfig struct {
	BatchSize int
	FlushInterval time.Duration
	MaxWorkers int
	QueueSize int
	BufferMaxSize int
}

type tcpDispatcher struct {
	logger *zap.Logger
	queue chan []byte
	config DispatcherConfig
	workerWG sync.WaitGroup
	bytePool *sync.Pool
}

func NewTCPDispatcher(logger *zap.Logger, cfg DispatcherConfig) domain.EventDispatcher {
	d := &tcpDispatcher{
		logger: logger,
		queue:  make(chan []byte, cfg.QueueSize),
		config: cfg,
		bytePool: &sync.Pool{
			New: func() interface{} {
				return make([]byte, cfg.BufferMaxSize)
			},
		},
	}

	d.start()
	return d
}

func (d *tcpDispatcher) Enqueue(ctx *fasthttp.RequestCtx, data []byte) error {
	if len(data) > d.config.BufferMaxSize {
        return domain.ErrPayloadTooLarge
    }

	buf := d.bytePool.Get().([]byte)
	
	n := copy(buf, data)
	item := buf[:n]

	select {
	case d.queue <- item:
		return nil
	default:
		d.bytePool.Put(buf)
		return domain.ErrDispatchQueueIsFull
	}
}

func (d *tcpDispatcher) start() {
	for i := 0; i < d.config.MaxWorkers; i++ {
		d.workerWG.Add(1)
		go d.worker()
	}
}

func (d *tcpDispatcher) worker() {
	defer d.workerWG.Done()
	
	batch := make([][]byte, 0, d.config.BatchSize)
	ticker := time.NewTicker(d.config.FlushInterval) 
    defer ticker.Stop()

	for {
		select {
		case item, ok := <-d.queue:
			if !ok {
				d.shutdown(batch)
				return
			}

			batch = append(batch, item)

			if len(batch) >= d.config.BatchSize {
				d.flush(batch)
				batch = batch[:0]
			}
		case <-ticker.C:
			if len(batch) > 0 {
				d.flush(batch)
				batch = batch[:0]
			}
		}
	}
}

func (d *tcpDispatcher) flush(batch [][]byte) {
	// Week 2: Implement length-prefixed TCP write to Rust Log Engine

	d.logger.Debug("flushing batch", zap.Int("count", len(batch)))

	for i := range batch {
		buf := batch[i]
		d.bytePool.Put(buf[:cap(buf)])
		batch[i] = nil 
	}
}

func (d *tcpDispatcher) shutdown(batch [][]byte) {
	if len(batch) > 0 {
		d.flush(batch)
	}
}

func (d *tcpDispatcher) Close() {
	d.logger.Info("shutting down dispatcher...")

	close(d.queue)

	d.workerWG.Wait()

	d.logger.Info("dispatcher shut down gracefully")
}