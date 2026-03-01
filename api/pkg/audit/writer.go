package audit

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/and-period/furumaru/api/pkg/log"
)

const (
	defaultBufferSize    = 1024
	defaultBatchSize     = 100
	defaultFlushInterval = 5 * time.Second
)

// Store はaudit logのバッチ書き込みに使用するインターフェース
type Store interface {
	BatchCreate(ctx context.Context, logs entity.AuditLogs) error
}

// Writer は非同期でaudit logをバッチ書き込みするライター
type Writer struct {
	store         Store
	ch            chan *entity.AuditLog
	batchSize     int
	flushInterval time.Duration
	wg            sync.WaitGroup
	done          chan struct{}
}

type WriterOption func(*Writer)

func WithBufferSize(size int) WriterOption {
	return func(w *Writer) {
		w.ch = make(chan *entity.AuditLog, size)
	}
}

func WithBatchSize(size int) WriterOption {
	return func(w *Writer) {
		w.batchSize = size
	}
}

func WithFlushInterval(d time.Duration) WriterOption {
	return func(w *Writer) {
		w.flushInterval = d
	}
}

// NewWriter は非同期audit logライターを作成し、バックグラウンドgoroutineを開始する
func NewWriter(store Store, opts ...WriterOption) *Writer {
	w := &Writer{
		store:         store,
		ch:            make(chan *entity.AuditLog, defaultBufferSize),
		batchSize:     defaultBatchSize,
		flushInterval: defaultFlushInterval,
		done:          make(chan struct{}),
	}
	for _, opt := range opts {
		opt(w)
	}

	w.wg.Add(1)

	go w.run()

	return w
}

// Send はaudit logをライターに送信する。チャネルが満杯の場合はドロップする。
func (w *Writer) Send(log *entity.AuditLog) {
	select {
	case w.ch <- log:
	default:
		slog.Warn("Audit log dropped: buffer full")
	}
}

// Close はライターを停止し、残りのログをフラッシュする
func (w *Writer) Close() {
	close(w.done)
	w.wg.Wait()
}

func (w *Writer) run() {
	defer w.wg.Done()

	ticker := time.NewTicker(w.flushInterval)
	defer ticker.Stop()

	batch := make(entity.AuditLogs, 0, w.batchSize)

	for {
		select {
		case entry := <-w.ch:
			batch = append(batch, entry)
			if len(batch) >= w.batchSize {
				w.flush(batch)
				batch = make(entity.AuditLogs, 0, w.batchSize)
			}
		case <-ticker.C:
			if len(batch) > 0 {
				w.flush(batch)
				batch = make(entity.AuditLogs, 0, w.batchSize)
			}
		case <-w.done:
			// ドレイン: チャネルに残っているログを全て取得
			for {
				select {
				case entry := <-w.ch:
					batch = append(batch, entry)
				default:
					if len(batch) > 0 {
						w.flush(batch)
					}

					return
				}
			}
		}
	}
}

func (w *Writer) flush(logs entity.AuditLogs) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := w.store.BatchCreate(ctx, logs)
	if err != nil {
		slog.Error("Failed to flush audit logs",
			log.Error(err),
			slog.Int("count", len(logs)),
		)
	}
}
