FROM golang:latest

WORKDIR /app/admin-api
COPY . ./
RUN go mod download
RUN go build -o /admin-api

EXPOSE 9092

CMD ["/admin-api"]
