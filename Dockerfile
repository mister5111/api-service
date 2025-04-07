FROM golang:1.24.0-alpine AS builder

RUN apk add --no-cache gcc musl-dev sqlite-dev

WORKDIR /api-service

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENV CGO_ENABLED=1
ENV GOOS=linux
ENV GOARCH=amd64

RUN go build -o api-service main.go
RUN chmod +x api-service

FROM alpine:latest

RUN apk add --no-cache sqlite-libs

WORKDIR /var/local/api-service

COPY --from=builder /api-service/api-service .
COPY --from=builder /api-service/conf ./conf
COPY --from=builder /api-service/storage/storage.db ./storage/storage.db

EXPOSE 80

CMD ["./api-service"]
