migration:
	@migrate create -ext sql -dir migrations $(filter-out $@,$(MAKECMDGOALS))

migration-up: migration-drop
	@migrate -database "postgres://postgres:password@localhost:5432/test_db?sslmode=disable" -path migrations up


migration-down:
	@migrate -database "postgres://postgres:password@localhost:5432/test_db?sslmode=disable" -path migrations down

migration-drop:
	@migrate -database "postgres://postgres:password@localhost:5432/test_db?sslmode=disable" -path migrations drop -f






