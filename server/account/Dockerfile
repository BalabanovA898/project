FROM golang:1.25.9-alpine

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o account ./src

EXPOSE 8080

CMD ["./account"]
