# Use Golang base image
FROM golang:1.22-alpine as builder

WORKDIR /app

RUN apk add --no-cache gcc musl-dev sqlite-dev



# Copy only the necessary files
COPY go.mod go.sum ./
RUN go mod tidy

COPY . .

# Build the app binary
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o app ./cmd/app

# Use a smaller image for running the app
FROM alpine:latest
WORKDIR /app

RUN apk add --no-cache sqlite-libs

COPY --from=builder /app/app /app/app

EXPOSE 8080
CMD ["./app"]