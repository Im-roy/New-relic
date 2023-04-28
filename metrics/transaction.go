package metrics

import (
	"context"

	"github.com/newrelic/go-agent/v3/newrelic"
)

type Transaction struct {
	App *newrelic.Application
}

func(t Transaction) CreateTransaction(ctx context.Context, transactionName string) *newrelic.Transaction {
	txn := t.App.StartTransaction(transactionName)
	return txn
}

func(t Transaction) UpdateContextWithTransaction(ctx context.Context, txn *newrelic.Transaction) context.Context {
	ctx = newrelic.NewContext(ctx, txn)
	return ctx
}

func(t Transaction) MarkError(ctx context.Context, err error) {
	txn := newrelic.FromContext(ctx)
	txn.NoticeError(err)
}

// func(t Transaction) InsertDistributedTraceHeaders(ctx context.Context, )