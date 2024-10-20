FROM golang:1.23 as golang

WORKDIR /app

COPY . .

RUN go mod download && go mod verify

RUN CGO_ENABLED=0 go build -v -o main

RUN ls

FROM alpine

COPY --from=golang /app/main .

RUN mkdir /etc/webhook

EXPOSE 8080

CMD ["/main"]
