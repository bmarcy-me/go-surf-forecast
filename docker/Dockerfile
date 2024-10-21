FROM golang:1.23-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o setup_db cmd/db/setup_db.go

RUN go build -o server cmd/server/main.go

CMD ["./server"]