FROM golang:alpine

RUN apk update && apk add git

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH


RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH"
WORKDIR $GOPATH/src/user-service

COPY . .

RUN go get -u github.com/golang/dep/cmd/dep
RUN dep init && dep ensure

RUN go build -o user-service

ENTRYPOINT ["./user-service"]

EXPOSE 3000