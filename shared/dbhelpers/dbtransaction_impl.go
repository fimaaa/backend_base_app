package dbhelpers

import (
    "context"
)

// WithoutTransaction is helper function that simplify the readonly db
func WithoutTransaction(ctx context.Context, trx WithoutTransactionDB, trxFunc func(dbCtx context.Context) error) error {
	dbCtx, err := trx.GetDatabase(ctx)
	if err != nil {
		return err
	}

	defer func(trx WithoutTransactionDB, ctx context.Context) {
		err := trx.Close(ctx)
		if err != nil {
			return
		}
	}(trx, dbCtx)

	return trxFunc(dbCtx)
}

// WithTransaction is helper function that simplify the transaction execution handling
func WithTransaction(ctx context.Context, trx WithTransactionDB, trxFunc func(dbCtx context.Context) error) error {
	dbCtx, err := trx.BeginTransaction(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			err = trx.RollbackTransaction(dbCtx)
			panic(p)

		} else if err != nil {
			err = trx.RollbackTransaction(dbCtx)

		} else {
			err = trx.CommitTransaction(dbCtx)

		}
	}()

	err = trxFunc(dbCtx)
	return err
}
