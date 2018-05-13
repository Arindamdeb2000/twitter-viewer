FROM golang:1.9.1

ENV GOPATH=/go:/go/src/github.com/morrah77/twitter-viewer
WORKDIR /go/src/github.com/morrah77/twitter-viewer
COPY . .

RUN go get -u github.com/golang/dep/cmd/dep

RUN dep ensure
RUN go build -o twitter-viewer main.go
CMD ./twitter-viewer
