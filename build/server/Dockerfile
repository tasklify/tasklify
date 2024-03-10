# Build
FROM golang:latest as build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /entrypoint ./cmd/server/main.go

# Deploy
FROM gcr.io/distroless/static-debian12:latest as release-stage

WORKDIR /

COPY --chown=nonroot --from=build-stage /entrypoint /entrypoint
COPY --chown=nonroot --from=build-stage /app/static /static

ENV PORT 8080
EXPOSE $PORT

USER nonroot:nonroot
ENTRYPOINT ["/entrypoint"]
