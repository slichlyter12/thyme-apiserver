FROM golang:alpine as builder
WORKDIR /app
ADD . /app
RUN go build -o apiserver .

FROM alpine:latest as runner
COPY --from=builder /app/apiserver /app/apiserver
EXPOSE 8080

CMD ["/app/apiserver"]
