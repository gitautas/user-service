FROM golang:alpine as builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o main src/main.go

FROM alpine

COPY --from=builder /app/main /main

CMD [ "/main" ]
