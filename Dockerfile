# Build
FROM golang:latest as build-stage

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /tasklify ./cmd/server/main.go

# Deploy
FROM gcr.io/distroless/static-debian12:latest as release-stage

COPY --from=build-stage /tasklify /tasklify
COPY --from=build-stage /app/static /static

ENV PORT ${ PORT }
EXPOSE $PORT

USER nonroot:nonroot
ENTRYPOINT ["/tasklify"]
