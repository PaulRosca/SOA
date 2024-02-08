FROM golang:1.21.6 as builder

WORKDIR /app
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o orders

FROM alpine:3.14

WORKDIR /app
COPY --from=builder /app/orders /app/

EXPOSE 6666
CMD ["/app/orders"]
