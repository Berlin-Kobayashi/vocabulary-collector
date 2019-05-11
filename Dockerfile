FROM golang:1.12 as build

WORKDIR /go/src/github.com/DanShu93/vocabulary-collector

COPY . .

RUN go get ./...

CMD go run app/app.go
