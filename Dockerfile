FROM golang:1.13 AS builder
WORKDIR /github.com/tjper/shoppingcart/server/
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o shoppingcart .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /github.com/tjper/shoppingcart/server/shoppingcart .

EXPOSE 8080
ENTRYPOINT ["./shoppingcart"]
CMD ["serve"]
