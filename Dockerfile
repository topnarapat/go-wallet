FROM golang:1.19-alpine as build-base

WORKDIR /app

COPY go.mod .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go test --tags=unit -v ./...

RUN go build -o ./out/go-wallet .

# ====================

FROM alpine:3.16.2

COPY --from=build-base /app/out/go-wallet /app/go-wallet

CMD ["/app/go-wallet"]