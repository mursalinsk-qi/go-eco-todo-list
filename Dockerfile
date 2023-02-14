FROM golang:1.19.0
WORKDIR /usr/src/app
COPY . .
RUN go mod download
RUN go build -o main .
CMD ["./main"]