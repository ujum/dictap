FROM golang:1.17.0-alpine3.14 AS builder

WORKDIR /home/app/dictup
COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY ./ ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./out/app ./cmd/dictup/main.go

FROM alpine:3.14

WORKDIR /home
COPY --from=builder /home/app/dictup/out/app ./
COPY --from=builder /home/app/dictup/configs/ ./configs/

EXPOSE 8080

ENTRYPOINT ["./app"]