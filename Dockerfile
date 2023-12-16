FROM golang:1.20.7

WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go build -o /build

ENV HOSTNAME "0.0.0.0"

EXPOSE 8080

CMD ["/build"]
