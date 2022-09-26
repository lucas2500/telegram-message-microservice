FROM golang:1.19-alpine AS api-builder
RUN apk add build-base
RUN mkdir /app
COPY . /app
WORKDIR /app/api
RUN go build server.go

FROM alpine AS api-service
RUN mkdir /app
RUN mkdir /api
WORKDIR /app
COPY --from=api-builder /app/*.env /app/
COPY --from=api-builder /app/api /app/api/
WORKDIR /app/api
EXPOSE 3001
CMD ["/app/api/server"]