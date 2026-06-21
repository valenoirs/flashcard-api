export APLICATION_NAME=flashcard
export APPLICATION_ENV=development
export APPLICATION_VERSION=1.0.0
export APPLICATION_PORT=5001

export DB_CONNECTION_STRING=postgres://lucky:lucky@localhost:5433/flashcard?sslmode=disable
export DB_SCHEMA_NAME=flashcard_schema
export DB_MAX_LIFETIME=300s
export DB_MAX_IDLE_TIME=300s
export DB_MAX_IDLE_CONNECTION=10
export DB_MAX_OPEN_CONNECTION=100
export DB_SLOW_QUERY_THRESHOLD=200ms

go run ./cmd/
