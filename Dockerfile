FROM golang:1.17.0-alpine3.14 AS builder

RUN apk add make

WORKDIR /home/app/dictup
COPY go.mod go.sum ./

RUN go mod download

COPY ./ ./

RUN make build

FROM scratch

WORKDIR /home/dictup/
COPY --from=builder /home/app/dictup/out/ ./

EXPOSE 8080

ENTRYPOINT ["./dictup"]