Edit from [sikozonpc](https://github.com/sikozonpc/ecom)

## Stack

- config: godotenv
- router: chi
- database: postgres, pgx
- orm: sqlc
- auth: jwt

adding transaction, database conning pool.

## get start


```bash
# 1
git clone ...

# 2 use your config for database and jwt
cp .env.example .env

# 3 edit database config in Makefile and run
make migrate-up

# 4. 
go run cmd/main.go
```
