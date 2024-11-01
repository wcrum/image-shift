FROM cgr.dev/chainguard/go AS golang

WORKDIR /app

COPY . .

RUN go mod download && go mod verify

RUN GOEXPERIMENT=boringcrypto CGO_ENABLED=0 go build -v -o main -ldflags="-s -w"

FROM alpine:latest

COPY --from=golang /app/main .

RUN mkdir /etc/webhook

EXPOSE 8080

CMD ["/main"]
