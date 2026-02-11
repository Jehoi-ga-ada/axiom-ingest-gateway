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
}

type tcpDispatcher struct {
	logger *zap.Logger
	queue chan []byte
	config DispatcherConfig
	workerWG sync.WaitGroup
}

func NewTCPDispatcher(logger *zap.Logger, cfg DispatcherConfig) domain.EventDispatcher {
	d := &tcpDispatcher{
		logger: logger,
		queue:  make(chan []byte, cfg.QueueSize),
		config: cfg,
	}

	d.start()
	return d
}

func (d *tcpDispatcher) Enqueue(ctx *fasthttp.RequestCtx, data []byte) error {
	item := make([]byte, len(data))
	copy(item, data)

	select {
	case d.queue <- item:
		return nil
	default:
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
	timer := time.NewTimer(d.config.FlushInterval)
	defer timer.Stop()

	for {
		select {
		case item, ok := <-d.queue:
			if !ok {
				d.flush(batch)
				return
			}
			
			batch = append(batch, item)

			if len(batch) >= d.config.BatchSize {
				d.flush(batch)
				batch = batch[:0]

				if !timer.Stop() {
					select {
					case <-timer.C:
					default:
					}
				}
				timer.Reset(d.config.FlushInterval)
			}
		case <-timer.C:
			if len(batch) > 0 {
				d.flush(batch)
				batch = batch[:0]
			}

			timer.Reset(d.config.FlushInterval)
		}
	}
}

func (d *tcpDispatcher) flush(batch [][]byte) {
	if len(batch) == 0 {
		return
	}

	// Week 2: Implement length-prefixed TCP write to Rust Log Engine

	d.logger.Debug("flushing batch to log engine", 
		zap.Int("count", len(batch)),
		zap.String("target", "axiom-event-log-engine"),
	)
}

func (d *tcpDispatcher) Close() {
	d.logger.Info("shutting down dispatcher...")

	close(d.queue)

	d.workerWG.Wait()

	d.logger.Info("dispatcher shut down gracefully")
}