FROM golang:1.18.2-alpine3.15
RUN mkdir /app
ADD . /app/
WORKDIR /app
RUN go mod download
RUN go build main.go
CMD ["/app/main"]
