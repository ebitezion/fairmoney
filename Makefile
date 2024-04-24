RUN:
	go run ./cmd/api

CREATE_BANK_TABLE_MIGRATION:
	migrate create -seq -ext=.sql -dir=./migrations fairmoney

EXECUTE_UP_MIGRATIONS:
	migrate -path=./migrations -database=$$DATABASE_DSN up

EXECUTE_DOWN_MIGRATIONS:
	migrate -path=./migrations -database=$$DATABASE_DSN down

