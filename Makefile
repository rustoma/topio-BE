include .env

migration_up:
	migrate -path database/migration/ -database "postgresql://${dbuser}@${host}:${dbport}/${dbname}?sslmode=disable" -verbose up

migration_down:
	migrate -path database/migration/ -database "postgresql://${dbuser}@${host}:${dbport}/${dbname}?sslmode=disable" -verbose down

migration_fix:
	migrate -path database/migration/ -database "postgresql://${dbuser}@${host}:${dbport}/${dbname}?sslmode=disable" force VERSION