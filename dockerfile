FROM golang:1.22-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o /strong_password ./app/main.go

EXPOSE 8080

CMD ["/strong_password"]