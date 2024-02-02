postgres:
	docker run --name postgres12 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=12345678 -d postgres:12-alpine
createdb:
	docker exec -it postgres12 createdb --username=root --owner=root simple_bank
dropdb:
	docker exec -it postgres12 dropdb simple_bank
migrateup:
	# 迁移并创建数据库
	migrate -path db/migration -database "postgres://root:12345678@localhost:5432/simple_bank?sslmode=disable" -verbose up
migratedown:
	migrate -path db/migration -database "postgres://root:12345678@localhost:5432/simple_bank?sslmode=disable" -verbose down
showall:
	ls -a
sqlc:
	sqlc generate
test:
	go test -v -cover ./...
.PHONY: createdb postgres dropdb migrateup migratedown showall sqlc test