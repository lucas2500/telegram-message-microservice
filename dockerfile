FROM golang:1.18.1-alpine AS CONSUMER
RUN mkdir /app
ADD . /app
WORKDIR /app
RUN apk add build-base
RUN go build main.go
CMD ["/app/main"]

FROM golang:1.18.1-alpine AS API
RUN mkdir /app
ADD . /app
WORKDIR /app
RUN apk add build-base
RUN go build api/main.go
EXPOSE 3001
CMD ["/app/api/main"]