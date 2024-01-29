package db

import (
	"context"
	"database/sql"
	"fmt"
)

// 做一个事务
// 合成，使store具有*Queries的所有功能
type Store struct {
	*Queries
	db *sql.DB
}

// db就是数据库
// Queries是数据库中的查询，具有上下文功能 由sqlc生成
func NewStore(db *sql.DB) *Store {
	return &Store{
		Queries: New(db),
		db:      db,
	}
}

// fn 回调函数
func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// 具有事务的查询
	q := New(tx)

	err = fn(q)
	// 说明需要回滚
	if err != nil {
		// 回滚失败和交易失败事务的错误
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err:%v, rb err:%v\n", err, rbErr)
		}
		// 只有交易失败的错误
		return err
	}

	return tx.Commit()
}

// 事务交易结构体
// 转账人ID
// 收款人ID
// 金额
type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

// 事务交易结果结构体
// 交易记录
type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"` //转账人账户更新
	ToAccount   Account  `json:"to_account"`   //收款人账户更新
	FromEntry   Entry    `json:"from_entry"`   //转账记录
	ToEntry     Entry    `json:"to_entry"`     //收款记录
}

func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}

		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		// TODO: update

		return nil
	})

	return result, err
}
