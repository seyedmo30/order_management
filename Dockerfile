FROM golang:1.22-alpine as builder
WORKDIR /app
RUN apk add  gcc musl-dev sqlite-dev
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o app ./cmd/app
FROM alpine:3.21
WORKDIR /app
RUN apk add  sqlite-libs
COPY --from=builder /app/app /app/app
EXPOSE 8080
CMD ["./app"]