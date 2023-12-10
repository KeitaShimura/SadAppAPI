FROM golang:1.20.7

WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .

ENV GOOS=linux GOARCH=amd64
RUN go build -o main ./main.go

EXPOSE 8080
CMD ["/app/main"]
