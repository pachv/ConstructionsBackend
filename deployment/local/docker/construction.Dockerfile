FROM golang:1.24-alpine3.22 AS builder

RUN apk --no-cache add ca-certificates git

WORKDIR /app


COPY services/constructions/go.mod services/constructions/go.sum ./
RUN go mod download

COPY services/constructions .

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build -o main ./cmd/main.go

FROM alpine:3.22

# Устанавливаем только сертификаты
RUN apk --no-cache add ca-certificates

RUN apk --no-cache add ca-certificates
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.edge.kernel.org/g' /etc/apk/repositories \
    && apk update \
    && apk add --no-cache tzdata

WORKDIR /app

COPY --from=builder /app/main .

EXPOSE 80

CMD ["./main"]