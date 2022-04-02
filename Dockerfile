FROM golang:1.18-alpine
WORKDIR /app
RUN apk update && apk add git
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -o main .
CMD ["./main"] 