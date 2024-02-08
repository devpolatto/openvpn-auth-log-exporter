FROM golang:latest as builder
WORKDIR /go/src/app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN go build -o main .
EXPOSE 9177
CMD ["./main"]
