FROM golang:1.23-alpine AS build
WORKDIR /app/services/attendance-service
COPY . /app
RUN apk add --no-cache ca-certificates && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server ./cmd/server

FROM gcr.io/distroless/base-debian12:nonroot
WORKDIR /app
COPY --from=build /app/services/attendance-service/server /app/server
EXPOSE 9093
USER nonroot:nonroot
ENTRYPOINT ["/app/server"]
