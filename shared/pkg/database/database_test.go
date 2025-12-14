package database

import (
	"context"
	"errors"
	"testing"

	"github.com/jackc/pgx/v5"
)

type fakeTx struct {
	committed bool
	rolled    bool
}

func (f *fakeTx) Commit(ctx context.Context) error {
	f.committed = true
	return nil
}
func (f *fakeTx) Rollback(ctx context.Context) error {
	f.rolled = true
	return nil
}

type fakeStarter struct {
	failBegin bool
	lastTx    *fakeTx
}

func (s *fakeStarter) BeginTx(ctx context.Context, txOptions pgx.TxOptions) (Tx, error) {
	if s.failBegin {
		return nil, errors.New("begin failed")
	}
	s.lastTx = &fakeTx{}
	return s.lastTx, nil
}

type fakeTxCommitErr struct {
	fakeTx
}

func (t *fakeTxCommitErr) Commit(ctx context.Context) error {
	return errors.New("commit failed")
}

type fakeStarterCommitErr struct{}

func (s *fakeStarterCommitErr) BeginTx(ctx context.Context, txOptions pgx.TxOptions) (Tx, error) {
	return &fakeTxCommitErr{}, nil
}

func TestWithTxCommit(t *testing.T) {
	s := &fakeStarter{}
	err := WithTx(context.Background(), s, func(tx Tx) error {
		return nil
	})
	if err != nil {
		t.Fatalf("WithTx error: %v", err)
	}
	if s.lastTx == nil || !s.lastTx.committed {
		t.Fatalf("expected commit to be called")
	}
}

func TestWithTxRollback(t *testing.T) {
	s := &fakeStarter{}
	err := WithTx(context.Background(), s, func(tx Tx) error {
		return errors.New("fail")
	})
	if err == nil {
		t.Fatalf("expected error")
	}
	if s.lastTx == nil || !s.lastTx.rolled {
		t.Fatalf("expected rollback to be called")
	}
}

func TestCommitErrorPropagates(t *testing.T) {
	err := WithTx(context.Background(), &fakeStarterCommitErr{}, func(tx Tx) error {
		return nil
	})
	if err == nil {
		t.Fatalf("expected commit error")
	}
}

func TestBeginTxFail(t *testing.T) {
	s := &fakeStarter{failBegin: true}
	err := WithTx(context.Background(), s, func(tx Tx) error {
		return nil
	})
	if err == nil {
		t.Fatalf("expected begin error")
	}
}

func TestConnectAttemptValidURL(t *testing.T) {
	// URL valid tetapi koneksi akan gagal karena tidak ada server Postgres lokal
	p, err := Connect(context.Background(), "postgres://user:pass@localhost:5432/dbname")
	if err != nil {
		// Terima kesalahan jika implementasi pool mencoba koneksi eager
		return
	}
	if p == nil {
		t.Fatalf("expected non-nil pool")
	}
	p.Close()
}

func TestConnectInvalidURL(t *testing.T) {
	_, err := Connect(context.Background(), "invalid-url")
	if err == nil {
		t.Fatalf("expected parse error")
	}
}
