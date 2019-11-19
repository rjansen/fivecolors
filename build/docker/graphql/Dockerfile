FROM golang:1.13.4 as builder
WORKDIR /fivecolors
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor -o graphql cmd/graphql/main.go

FROM alpine:3.10
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /fivecolors/graphql .
CMD ["./graphql"]
