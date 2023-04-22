FROM library/golang:1.20-alpine as builder
WORKDIR /app

RUN apk add --no-cache ca-certificates && \
adduser -u 10001 --disabled-password app

COPY go.mod go.sum ./
RUN go mod download

COPY ./ ./
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-w -s" -o /app/main .

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
USER app

COPY --from=builder /app/main .

EXPOSE 8080

ENTRYPOINT ["/main"]
