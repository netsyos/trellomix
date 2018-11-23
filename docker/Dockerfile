FROM golang:1.8

WORKDIR /go/src/trellomix
COPY ./src/ .
RUN ls /go/src/trellomix

RUN go get -d -v ./...
RUN go install -v ./...

CMD ["trellomix"]