FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY . .
RUN go build -o kvStorage main.go
EXPOSE 8081
CMD ["./kvStorage"]
