FROM golang:1.24.0 AS builder

WORKDIR /api-service

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o api-service main.go
RUN chmod +x api-service

FROM scratch

WORKDIR /var/www/api-service
COPY --from=builder /api-service/api-service .
COPY --from=builder /api-service/conf ./conf

EXPOSE 80

CMD ["./api-service"]