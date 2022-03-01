FROM golang:1.17-alpine
RUN mkdir /app
RUN apk add build-base
LABEL authors="Li-Khan & quazar"
ADD . /app/
WORKDIR /app
RUN go build -o main ./cmd/
CMD ["/app/main"]