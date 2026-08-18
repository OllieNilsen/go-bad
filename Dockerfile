FROM golang:latest

WORKDIR /app
COPY . .
RUN go build -o badservice .

EXPOSE 8080
CMD ["./badservice"]
