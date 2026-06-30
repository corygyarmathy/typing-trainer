module github.com/corygyarmathy/typist

go 1.26.4

// Dependencies will be added as each phase is implemented
// Expected core set after phase 4:
//   - github.com/go-chi/chi/v5
//   - github.com/jackc/pgx/v5
//   - github.com/pressly/goose/v3
//   - github.com/oapi-codegen/runtime
//   - github.com/golang-jwt/jwt/v5
//   - github.com/prometheus/client_golang
//   - golang.org/x/crypto (for bcrypt/argon2id)

require (
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/pgx/v5 v5.10.0 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	golang.org/x/sync v0.17.0 // indirect
	golang.org/x/text v0.29.0 // indirect
)
