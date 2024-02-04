package main

import (
    "database/sql"
    "log"

    _ "github.com/lib/pq"
    "techschool/simplebank/api"
    db "techschool/simplebank/db/sqlc"
    "techschool/simplebank/util"
)

func main() {

    // 从环境变量中加载配置 .表示当前目录
    config, err := util.LoadConfig(".")
    if err != nil {
        log.Fatal("cannot load config:", err)
    }

    conn, err := sql.Open(config.DBDriver, config.DBSource)
    if err != nil {
        log.Fatal("cannot connect to db:", err)
    }
    store := db.NewStore(conn)
    server := api.NewServer(store)

    err = server.Start(config.ServerAddress)
    if err != nil {
        log.Fatal("cannot start server:", err)
    }
}
