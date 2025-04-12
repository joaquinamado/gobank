# Dockerfile
FROM golang:1.23.4

# Install air
RUN go install github.com/air-verse/air@latest

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

EXPOSE 8080

CMD ["air"]

