FROM golang:1.12 as build

WORKDIR /go/src/github.com/DanShu93/vocabulary-collector

COPY . .

RUN GO111MODULE=on go get ./...

CMD GO111MODULE=on go run app/app.go
