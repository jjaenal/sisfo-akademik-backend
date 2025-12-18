FROM golang:1.25-alpine AS build
WORKDIR /app/services/admission-service
COPY . /app
RUN apk add --no-cache ca-certificates && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server ./cmd/server

FROM gcr.io/distroless/base-debian12:nonroot
WORKDIR /app
COPY --from=build /app/services/admission-service/server /app/server
EXPOSE 9095
USER nonroot:nonroot
ENTRYPOINT ["/app/server"]
