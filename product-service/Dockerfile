FROM golang:alpine

RUN apk update && apk add git

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH


RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH"
WORKDIR $GOPATH/src/product-service

COPY . .

RUN go get -u github.com/golang/dep/cmd/dep
RUN dep init && dep ensure -v

RUN go build -o product-service

ENTRYPOINT ["./product-service"]

EXPOSE 3001