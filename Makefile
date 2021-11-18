POSTGRESQL_URL="postgres://admin:admin@localhost:5432/conduit?sslmode=disable"

run:
	go run main.go

create-migration:
	migrate create -ext sql -dir postgres/migrations -seq $(file)

run-migration:
	migrate -database $(POSTGRESQL_URL) -path postgres/migrations up

down-migration:
	migrate -database $(POSTGRESQL_URL) -path postgres/migrations down $(v)

force-migration:
	migrate -database $(POSTGRESQL_URL) -path postgres/migrations force $(version)