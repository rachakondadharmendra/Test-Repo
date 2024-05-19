module backend_golang

go 1.22.1

require (
	github.com/golang/snappy v0.0.4 // indirect
	github.com/gorilla/mux v1.8.1 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/klauspost/compress v1.17.7 // indirect
	github.com/montanaflynn/stats v0.7.1 // indirect
	github.com/rs/cors v1.10.1 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.1.2 // indirect
	github.com/xdg-go/stringprep v1.0.4 // indirect
	github.com/youmark/pkcs8 v0.0.0-20201027041543-1326539a0a0a // indirect
	go.mongodb.org/mongo-driver v1.14.0 // indirect
	golang.org/x/crypto v0.21.0 // indirect
	golang.org/x/sync v0.6.0 // indirect
	golang.org/x/text v0.14.0 // indirect
)

replace (
    backend_golang/src/db => ./src/db
    backend_golang/src/models => ./src/models
    backend_golang/src/logger => ./src/logger
	backend_golang/src/handlers => ./src/handlers
)
