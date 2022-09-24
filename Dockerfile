FROM golang:1.19-alpine
RUN mkdir /app
ADD . /app
WORKDIR /app/api
RUN apk add build-base
RUN go build server.go
EXPOSE 3001
CMD ["/app/api/server"]