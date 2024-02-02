package db

import (
    "context"
    "fmt"
    "testing"

    "github.com/stretchr/testify/require"
)

func TestTransfertx(t *testing.T) {
    store := NewStore(testDb)

    // 测试账户1向账户2转账10
    // 1. 创建两个账户
    account1 := createRandomAccount(t)
    account2 := createRandomAccount(t)

    fmt.Println("before transfer:", account1.Balance, account2.Balance)

    // 2.使用5个goroutine并发执行转账操作
    n := 10
    amount := int64(10)

    // 3.接收错误和结果的通道
    errc := make(chan error)

    for i := 0; i < n; i++ {

        fromAccount := account1
        toAccount := account2

        if i%2 == 1 {
            fromAccount = account2
            toAccount = account1
        }

        go func() {
            _, err := store.TransferTx(context.Background(), TransferTxParams{
                FromAccountID: fromAccount.ID,
                ToAccountID:   toAccount.ID,
                Amount:        amount,
            })
            // 不能在这里使用require.NoError(t, err) 因为这里是在goroutine中,应该将错误返回到主测试函数，然后在那里处理错误
            errc <- err
        }()
    }
    // existed := map[int]bool{}
    // 4.检查错误
    for i := 0; i < n; i++ {
        err := <-errc
        require.NoError(t, err)

        // // 5.检查结果
        // res := <-resc
        // require.NotEmpty(t, res)
        //
        // transfer := res.Transfer
        // require.Equal(t, account1.ID, transfer.FromAccountID)
        // require.Equal(t, account2.ID, transfer.ToAccountID)
        // // 转账的金额
        // require.Equal(t, amount, transfer.Amount)
        // require.NotZero(t, transfer.ID)
        // require.NotZero(t, transfer.CreatedAt)
        //
        // // 6.检查转帐记录存在
        // _, err = store.GetTransfer(context.Background(), transfer.ID)
        // require.NoError(t, err)
        //
        // // 7.检查
        // fromEntry := res.FromEntry
        // require.NotEmpty(t, fromEntry)
        // require.Equal(t, account1.ID, fromEntry.AccountID)
        // require.Equal(t, -amount, fromEntry.Amount)
        // require.NotZero(t, fromEntry.ID)
        // require.NotZero(t, fromEntry.CreatedAt)
        // _, err = store.GetEntry(context.Background(), fromEntry.ID)
        // require.NoError(t, err)
        //
        // // 8.检查
        // toEntry := res.ToEntry
        // require.NotEmpty(t, toEntry)
        // require.Equal(t, account2.ID, toEntry.AccountID)
        // require.Equal(t, amount, toEntry.Amount)
        // require.NotZero(t, toEntry.ID)
        // require.NotZero(t, toEntry.CreatedAt)
        // _, err = store.GetEntry(context.Background(), toEntry.ID)
        // require.NoError(t, err)
        //
        // fromAccount := res.FromAccount
        // require.NotEmpty(t, fromAccount)
        // require.Equal(t, account1.ID, fromAccount.ID)
        //
        // toAccount := res.ToAccount
        // require.NotEmpty(t, toAccount)
        // require.Equal(t, account2.ID, toAccount.ID)
        //
        // fmt.Println("tx:", fromAccount.Balance, toAccount.Balance)
        //
        // // todo: 检查账户余额是否正确
        // // 原来的余额 和 转账后的余额做差值
        // diff1 := account1.Balance - fromAccount.Balance
        // diff2 := toAccount.Balance - account2.Balance
        // require.Equal(t, diff1, diff2)
        // // 差额是正数
        // require.True(t, diff1 > 0)
        // // 差额是转账金额的整数倍
        // require.True(t, diff1%amount == 0)
        //
        // k := int(diff1 / amount)
        // require.True(t, k >= 1 && k <= n)
        // require.NotContains(t, existed, k)
        // existed[k] = true
    }

    // 检查更新后的账户余额
    updateAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
    require.NoError(t, err)

    updateAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
    require.NoError(t, err)
    fmt.Println("after transfer:", updateAccount1.Balance, updateAccount2.Balance)
    require.Equal(t, account1.Balance, updateAccount1.Balance)
    require.Equal(t, account2.Balance, updateAccount2.Balance)

}
