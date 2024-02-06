FROM golang:1.21.2 as builder

WORKDIR /app
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o catalog

FROM alpine:3.14

WORKDIR /app
COPY --from=builder /app/catalog /app/

EXPOSE 5555
CMD ["/app/catalog"]
