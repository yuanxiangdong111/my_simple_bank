package db

import (
    "context"
    "testing"
    "time"

    "github.com/stretchr/testify/require"
    "techschool/simplebank/util"
)

func createRandomUser(t *testing.T) User {
    hashedPassword, err := util.HashPassword(util.RandomString(6))
    require.NoError(t, err)

    arg := CreateUserParams{
        Username:       util.RandomOwner(),
        HashedPassword: hashedPassword,
        FullName:       util.RandomOwner(),
        Email:          util.RandomEmail(),
    }

    user, err := testQueries.CreateUser(context.Background(), arg)
    // 测试没有错误
    require.NoError(t, err)
    // 测试非空
    require.NotEmpty(t, user)

    // 测试参数相等
    require.Equal(t, arg.Username, user.Username)
    require.Equal(t, arg.HashedPassword, user.HashedPassword)
    require.Equal(t, arg.FullName, user.FullName)
    require.Equal(t, arg.Email, user.Email)

    // 测试id和创建时间是自动生成的
    require.NotZero(t, user.CreatedAt)
    // 首次更改密码时间为空
    require.True(t, user.PasswordChangedAt.IsZero())
    return user
}

func TestCreateUser(t *testing.T) {
    createRandomUser(t)
}

func TestGetUser(t *testing.T) {
    user1 := createRandomUser(t)
    user2, err := testQueries.GetUser(context.Background(), user1.Username)
    require.NoError(t, err)
    require.NotEmpty(t, user2)

    // 参数是否相同
    require.Equal(t, user1.Username, user2.Username)
    require.Equal(t, user1.HashedPassword, user2.HashedPassword)
    require.Equal(t, user1.FullName, user2.FullName)
    require.Equal(t, user1.Email, user2.Email)

    // 时间戳相差
    require.WithinDuration(t, user1.PasswordChangedAt, user2.PasswordChangedAt, time.Second)
    require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
}

// func TestUpdateUser(t *testing.T) {
//     user1 := createRandomAccount(t)
//     arg := UpdateAccountParams{
//         ID:      user1.ID,
//         Balance: util.RandomMoney(),
//     }
//     user2, err := testQueries.UpdateAccount(context.Background(), arg)
//     require.NoError(t, err)
//     require.NotEmpty(t, user2)
//
//     // 参数是否相同
//     require.Equal(t, user1.ID, user2.ID)
//     require.Equal(t, user1.Owner, user2.Owner)
//     require.Equal(t, arg.Balance, user2.Balance)
//     require.Equal(t, user1.Currency, user2.Currency)
//
//     // 时间戳相差
//     require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
// }
//
// func TestDeleteAccount(t *testing.T) {
//     account1 := createRandomAccount(t)
//     err := testQueries.DeleteAccount(context.Background(), account1.ID)
//     require.NoError(t, err)
//
//     account2, err := testQueries.GetAccount(context.Background(), account1.ID)
//     require.Error(t, err)
//     require.EqualError(t, err, sql.ErrNoRows.Error())
//     require.Empty(t, account2)
// }
//
// func TestListAccounts(t *testing.T) {
//     for i := 0; i < 10; i++ {
//         createRandomAccount(t)
//     }
//
//     // 拿5-10位置的元素
//     arg := ListAccountsParams{
//         Limit:  5,
//         Offset: 5,
//     }
//
//     accounts, err := testQueries.ListAccounts(context.Background(), arg)
//     require.NoError(t, err)
//     require.Len(t, accounts, 5)
//
//     for _, account := range accounts {
//         require.NotEmpty(t, account)
//     }
// }
