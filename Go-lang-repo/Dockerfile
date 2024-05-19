FROM golang:1.22.1-alpine AS build
WORKDIR /server_build
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .

# Builds the application as a staticly linked one, to allow it to run on alpine
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Moving the binary to the 'final Image' to make it smaller
FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /server
RUN mkdir logs && touch server.log
COPY --from=build /server_build/main /server_build/.env /server_build/api_all.html /server_build/api_error.html /server/
EXPOSE 8080
CMD ["./main"]