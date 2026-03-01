package audit

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/and-period/furumaru/api/internal/user/entity"
	"github.com/stretchr/testify/assert"
)

type mockStore struct {
	mu      sync.Mutex
	logs    entity.AuditLogs
	callErr error
}

func (m *mockStore) BatchCreate(_ context.Context, logs entity.AuditLogs) error {
	if m.callErr != nil {
		return m.callErr
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	m.logs = append(m.logs, logs...)

	return nil
}

func (m *mockStore) getLogs() entity.AuditLogs {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.logs
}

func TestWriter_Send(t *testing.T) {
	t.Parallel()

	store := &mockStore{}
	w := NewWriter(store,
		WithBufferSize(10),
		WithBatchSize(2),
		WithFlushInterval(50*time.Millisecond),
	)

	log1 := &entity.AuditLog{ID: "1", AdminID: "admin1"}
	log2 := &entity.AuditLog{ID: "2", AdminID: "admin2"}

	w.Send(log1)
	w.Send(log2)

	// バッチサイズに到達したのでフラッシュされるはず
	time.Sleep(100 * time.Millisecond)
	assert.Len(t, store.getLogs(), 2)

	w.Close()
}

func TestWriter_FlushInterval(t *testing.T) {
	t.Parallel()

	store := &mockStore{}
	w := NewWriter(store,
		WithBufferSize(10),
		WithBatchSize(100), // 大きいバッチサイズ
		WithFlushInterval(50*time.Millisecond),
	)

	log1 := &entity.AuditLog{ID: "1", AdminID: "admin1"}
	w.Send(log1)

	// フラッシュ間隔で書き込まれるはず
	time.Sleep(150 * time.Millisecond)
	assert.Len(t, store.getLogs(), 1)

	w.Close()
}

func TestWriter_Close(t *testing.T) {
	t.Parallel()

	store := &mockStore{}
	w := NewWriter(store,
		WithBufferSize(10),
		WithBatchSize(100), // 大きいバッチサイズ
		WithFlushInterval(10*time.Second),
	)

	log1 := &entity.AuditLog{ID: "1", AdminID: "admin1"}
	w.Send(log1)

	// Close時に残りのログがフラッシュされるはず
	w.Close()
	assert.Len(t, store.getLogs(), 1)
}

func TestWriter_BufferFull(t *testing.T) {
	t.Parallel()

	store := &mockStore{}
	w := NewWriter(store,
		WithBufferSize(1),
		WithBatchSize(100),
		WithFlushInterval(10*time.Second),
	)

	// チャネルを埋める
	w.Send(&entity.AuditLog{ID: "1"})

	// 2つ目はドロップされるはず（ブロックしない）
	done := make(chan struct{})

	go func() {
		w.Send(&entity.AuditLog{ID: "2"})
		close(done)
	}()

	select {
	case <-done:
		// ブロックしなかったのでOK
	case <-time.After(1 * time.Second):
		t.Fatal("Send blocked on full buffer")
	}

	w.Close()
}
