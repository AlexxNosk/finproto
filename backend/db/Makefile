postgres:
	docker run --name finProto \
	-e POSTGRES_USER=root \
	-e POSTGRES_PASSWORD=root \
	-e POSTGRES_DB=finProto_db \
	-v pgdata:/var/lib/postgresql/data \
	-p 5434:5432 \
	-d postgres:17.5-bookworm	

createdb:
	docker exec -it finProto createdb --username=root --owner=root finProto_db

dropdb:
	docker exec -it finProto dropdb finProto_db

migrateup:
	migrate -path migration -database \
	"postgresql://root:root@localhost:5434/finProto_db?sslmode=disable" \
	-verbose up

migratedown:
	migrate -path migration \
	-database "postgresql://root:root@localhost:5434/finProto_db?sslmode=disable" \
	-verbose down

sqlc:
	sqlc generate

.PHONY: postgres createdb dropdb migrateup migratedown sqlc


# docker ps
# docker stop stock
# docker rm stock

