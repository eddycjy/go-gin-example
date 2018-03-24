FROM golang:latest

WORKDIR $GOPATH/src/github.com/EDDYCJY/go-gin-example
COPY . $GOPATH/src/github.com/EDDYCJY/go-gin-example
RUN go build .

EXPOSE 8000
ENTRYPOINT ["./go-gin-example"]

#FROM scratch

#WORKDIR $GOPATH/src/github.com/EDDYCJY/go-gin-example
#COPY . $GOPATH/src/github.com/EDDYCJY/go-gin-example

#EXPOSE 8000
#CMD ["./go-gin-example"]