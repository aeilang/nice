migration:
	@migrate create -ext sql -dir migrations $(filter-out $@,$(MAKECMDGOALS))

migration-up: migration-drop
	@migrate -database "postgres://postgres:password@localhost:5432/test_db?sslmode=disable" -path migrations up


migration-down:
	@migrate -database "postgres://postgres:password@localhost:5432/test_db?sslmode=disable" -path migrations down

migration-drop:
	@migrate -database "postgres://postgres:password@localhost:5432/test_db?sslmode=disable" -path migrations drop -f


postgres:
	docker run -d --name postgres \
		-e POSTGRES_PASSWORD=password \
		-e POSTGRES_DB=test_db \
		-v /home/lang/postgres/data:/var/lib/postgresql/data \
		-p 5432:5432 \
		postgres:alpine


redis:
	docker run -d --name redis -p 6379:6379 redis:alpine

