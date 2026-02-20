package infrastructure

import (
	"encoding/binary"
	"net"
	"sync"
	"time"

	"github.com/Jehoi-ga-ada/axiom-ingest-gateway/internal/features/ingest/domain"
	"go.uber.org/zap"
)

type DispatcherConfig struct {
	BatchSize     int
	FlushInterval time.Duration
	MaxWorkers    int
	MaxSenders int
	QueueSize     int
	BufferMaxSize int
	TargetAddr    string
	WriteTimeout  time.Duration
}

type tcpDispatcher struct {
	logger   *zap.Logger
	config   DispatcherConfig
	queue    chan []byte
	outbound chan [][]byte
	bytePool *sync.Pool
	batchPool *sync.Pool 
	workerWG sync.WaitGroup
	senderWG sync.WaitGroup 
	stop     chan struct{}
}

func NewTCPDispatcher(logger *zap.Logger, cfg DispatcherConfig) EventDispatcher {
	d := &tcpDispatcher{
		logger:   logger,
		config:   cfg,
		queue:    make(chan []byte, cfg.QueueSize),
		outbound: make(chan [][]byte, cfg.MaxWorkers),
		stop:     make(chan struct{}),
		bytePool: &sync.Pool{
			New: func() interface{} {
				return make([]byte, cfg.BufferMaxSize)
			},
		},
		batchPool: &sync.Pool{
			New: func() interface{} {
				return make([][]byte, 0, cfg.BatchSize)
			},
		},
	}

	for i := 0; i < cfg.MaxSenders; i ++ {
		d.senderWG.Add(1)
		go d.sender()
	}

	for i := 0; i < cfg.MaxWorkers; i++ {
		d.workerWG.Add(1)
		go d.worker()
	}

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
		d.bytePool.Put(buf[:cap(buf)])
		return domain.ErrDispatchQueueIsFull
	}
}

func (d *tcpDispatcher) worker() {
	defer d.workerWG.Done()

	batch := d.batchPool.Get().([][]byte)
	batch = batch[:0]
	ticker := time.NewTicker(d.config.FlushInterval)
	defer ticker.Stop()

	for {
		select {
		case item, ok := <-d.queue:
			if !ok {
				if len(batch) > 0 {
					d.outbound <- batch
				} else {
					d.batchPool.Put(batch)
				}
				return
			}

			batch = append(batch, item)

			if len(batch) >= d.config.BatchSize {
				d.outbound <- batch
				batch = d.batchPool.Get().([][]byte)
			}
		case <-ticker.C:
			if len(batch) > 0 {
				d.outbound <- batch
				batch = d.batchPool.Get().([][]byte)
			}
		case <-d.stop:
			if len(batch) > 0 {
				d.outbound <- batch
			} else {
				d.batchPool.Put(batch)
			}
			return
		}
	}
}

func (d *tcpDispatcher) sender() {
	defer d.senderWG.Done()
	var conn net.Conn
	ackBuf := make([]byte, 1)

	headerPool := make([]byte, d.config.BatchSize*4)

	cleanup := func() {
		if conn != nil {
			conn.Close()
			conn = nil
		}
	}

	for batch := range d.outbound {
		if conn == nil {
			var err error
			conn, err = net.DialTimeout("tcp", d.config.TargetAddr, 3*time.Second)
			if err != nil {
				d.logger.Error("dial failed", zap.Error(err))
				d.releaseBatch(batch)
				time.Sleep(500 * time.Millisecond) 
				continue
			}
			if tcpConn, ok := conn.(*net.TCPConn); ok {
				tcpConn.SetNoDelay(true)
			}
		}

		conn.SetWriteDeadline(time.Now().Add(d.config.WriteTimeout))
		
		vBuffers := make(net.Buffers, 0, len(batch)*2)
        
        for i, msg := range batch {
            start := i * 4
			header := headerPool[start : start+4]
            binary.BigEndian.PutUint32(header, uint32(len(msg)))
            vBuffers = append(vBuffers, header, msg)
        }

		_, err := vBuffers.WriteTo(conn)

		if err == nil {
            conn.SetReadDeadline(time.Now().Add(d.config.WriteTimeout))
            _, err = conn.Read(ackBuf)
        }

        if err != nil {
            d.logger.Error("transport error", zap.Error(err))
            cleanup()
        }
        
        d.releaseBatch(batch)
	}
	cleanup()
}

func (d *tcpDispatcher) releaseBatch(batch [][]byte) {
	for i := range batch {
		if batch[i] != nil {
			d.bytePool.Put(batch[i][:cap(batch[i])])
			batch[i] = nil
		}
	}

	d.batchPool.Put(batch[:0])
}

func (d *tcpDispatcher) Close() {
	d.logger.Info("shutting down dispatcher")
	
	close(d.stop)
	close(d.queue)
	d.workerWG.Wait()
	close(d.outbound)
	d.senderWG.Wait()
	
	d.logger.Info("dispatcher closed gracefully")
}