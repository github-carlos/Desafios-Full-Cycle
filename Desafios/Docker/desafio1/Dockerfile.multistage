FROM golang:1.18-alpine as builder

WORKDIR /app

RUN go mod init module
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /usr/local/bin/app ./...

CMD ["app"]

FROM scratch
WORKDIR /app
COPY --from=builder /usr/local/bin/app ./
CMD [ "./app" ]