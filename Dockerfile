FROM golang:alpine

WORKDIR /app

ADD . /app

RUN go build -o apiserver .

EXPOSE 8080

CMD ["/app/apiserver"]
