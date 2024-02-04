package api

import (
    "bytes"
    "database/sql"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/golang/mock/gomock"
    "github.com/stretchr/testify/require"
    mockdb "techschool/simplebank/db/mock"
    db "techschool/simplebank/db/sqlc"
    "techschool/simplebank/util"
)

func TestGetAccount(t *testing.T) {
    account := randomAccount()

    testCases := []struct {
        name          string
        accountID     int64
        buildStub     func(store *mockdb.MockStore) // 预期返回的数据
        checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
    }{
        {
            name:      "OK",
            accountID: account.ID,
            buildStub: func(store *mockdb.MockStore) {
                store.EXPECT().
                    GetAccount(gomock.Any(), gomock.Eq(account.ID)).
                    Times(1).
                    Return(account, nil)
            },
            checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
                require.Equal(t, http.StatusOK, recorder.Code)
                requireBodyMatchAccount(t, recorder.Body, account)
            },
        },
        {
            name:      "NotFount",
            accountID: account.ID,
            buildStub: func(store *mockdb.MockStore) {
                store.EXPECT().
                    GetAccount(gomock.Any(), gomock.Eq(account.ID)).
                    Times(1).
                    Return(db.Account{}, sql.ErrNoRows)
            },
            checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
                require.Equal(t, http.StatusNotFound, recorder.Code)
            },
        },
        {
            name:      "StatusInternalServerError",
            accountID: account.ID,
            buildStub: func(store *mockdb.MockStore) {
                store.EXPECT().
                    GetAccount(gomock.Any(), gomock.Eq(account.ID)).
                    Times(1).
                    Return(db.Account{}, sql.ErrConnDone)
            },
            checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
                require.Equal(t, http.StatusInternalServerError, recorder.Code)
            },
        },
        {
            name:      "InvalidID",
            accountID: 0,
            buildStub: func(store *mockdb.MockStore) {
                store.EXPECT().
                    GetAccount(gomock.Any(), gomock.Any()).
                    Times(0)
            },
            checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
                require.Equal(t, http.StatusBadRequest, recorder.Code)
            },
        },
    }

    for i := range testCases {
        tc := testCases[i]
        t.Run(tc.name, func(t *testing.T) {
            // 创建一个mock controller
            controller := gomock.NewController(t)
            defer controller.Finish()

            // 创建一个mock store
            store := mockdb.NewMockStore(controller)
            // 调用buildStub 功能是为了在mock store上设置预期行为
            tc.buildStub(store)

            // 创建一个server
            server := NewServer(store)
            // 创建一个http recorder 用于记录http response
            recorder := httptest.NewRecorder()
            url := fmt.Sprintf("/accounts/%d", tc.accountID)
            request, err := http.NewRequest(http.MethodGet, url, nil)
            require.NoError(t, err)

            // 调用server.router.ServeHTTP方法
            // 传入recorder和request
            // 对比预期结果
            server.router.ServeHTTP(recorder, request)
            tc.checkResponse(t, recorder)
        })
    }

}

func randomAccount() db.Account {
    return db.Account{
        ID:       util.RandomInt(1, 1000),
        Owner:    util.RandomOwner(),
        Balance:  util.RandomMoney(),
        Currency: util.RandomCurrency(),
    }
}

func requireBodyMatchAccount(t *testing.T, body *bytes.Buffer, account db.Account) {
    data, err := io.ReadAll(body)
    require.NoError(t, err)
    var gotAccount db.Account
    err = json.Unmarshal(data, &gotAccount)
    require.NoError(t, err)
    require.Equal(t, account, gotAccount)
}
