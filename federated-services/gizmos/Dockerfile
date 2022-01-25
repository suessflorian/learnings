FROM golang

WORKDIR /go/src/app

COPY . .

RUN go mod tidy

CMD go run .
