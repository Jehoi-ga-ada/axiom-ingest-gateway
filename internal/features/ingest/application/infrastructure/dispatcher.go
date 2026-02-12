package infrastructure

import (
	"bufio"
	"encoding/binary"
	"net"
	"sync"
	"time"

	"github.com/Jehoi-ga-ada/axiom-ingest-gateway/internal/features/ingest/domain"
	"go.uber.org/zap"
)

type DispatcherConfig struct {
	BatchSize int
	FlushInterval time.Duration
	MaxWorkers int
	QueueSize int
	BufferMaxSize int
	TargetAddr string
	WriteTimeout time.Duration
}

type tcpDispatcher struct {
	logger *zap.Logger
	queue chan []byte
	config DispatcherConfig
	workerWG sync.WaitGroup
	bytePool *sync.Pool

	mu sync.Mutex
	conn net.Conn
	bw *bufio.Writer
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

func (d *tcpDispatcher) Enqueue(data []byte) error {
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
				d.finalFlush(batch)
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
	d.mu.Lock()
	defer d.mu.Unlock()

	if err := d.ensureConnection(); err != nil {
		d.logger.Error("network failure", zap.Error(err))
		d.dropBatch(batch)
		return
	}

	d.conn.SetWriteDeadline(time.Now().Add(d.config.WriteTimeout))

	header := make([]byte, 4)
	for _, msg := range batch {
		binary.BigEndian.PutUint32(header, uint32(len(msg)))

		d.bw.Write(header)
		d.bw.Write(msg)
	}

	if err := d.bw.Flush(); err != nil {
		d.logger.Error("flush failed", zap.Error(err))
		d.closeConnection()
		d.dropBatch(batch)
		return
	}

	d.releaseBatch(batch)
}

func (d *tcpDispatcher) ensureConnection() error {
	if d.conn != nil {
		return nil
	}

	conn, err := net.DialTimeout("tcp", d.config.TargetAddr, 5 * time.Second)
	if err != nil {
		return err
	}

	if tcpConn, ok := conn.(*net.TCPConn); ok {
		tcpConn.SetNoDelay(true)
		tcpConn.SetKeepAlive(true)
	}

	d.conn = conn
	d.bw = bufio.NewWriterSize(conn, 128 * 1024)
	return nil
}

func (d *tcpDispatcher) closeConnection() {
	if d.conn != nil {
		d.conn.Close()
		d.conn = nil
		d.bw = nil
	}
}

func (d *tcpDispatcher) releaseBatch(batch [][]byte) {
	for i := range batch {
		d.bytePool.Put(batch[i][:cap(batch[i])])
		batch[i] = nil
	}
}

func (d *tcpDispatcher) dropBatch(batch [][]byte) {
	d.logger.Warn("dropping batch due to network error", zap.Int("size", len(batch)))
	d.releaseBatch(batch)
}

func (d *tcpDispatcher) finalFlush(batch [][]byte) {
	if len(batch) > 0 {
		d.flush(batch)
	}
}

func (d *tcpDispatcher) Close() {
	d.logger.Info("shutting down dispatcher...")

	close(d.queue)

	d.workerWG.Wait()

	d.mu.Lock()

	d.closeConnection()

	d.mu.Unlock()

	d.logger.Info("dispatcher shut down gracefully")
}