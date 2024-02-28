FROM golang:1.21 as builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 go build -o /server main.go

FROM gcr.io/distroless/base-debian11 as final

COPY --from=builder /server /server

ENV PORT 8080
EXPOSE $PORT

ENTRYPOINT ["/server"]

# # Build.
# FROM golang:1.20 AS build-stage
# WORKDIR /app
# COPY go.mod go.sum ./
# RUN go mod download
# COPY . /app
# RUN CGO_ENABLED=0 GOOS=linux go build -o /entrypoint
# 
# # Deploy.
# FROM gcr.io/distroless/static-debian11 AS release-stage
# WORKDIR /
# COPY --from=build-stage /entrypoint /entrypoint
# COPY --from=build-stage /app/assets /assets
# EXPOSE 8080
# USER nonroot:nonroot
# ENTRYPOINT ["/entrypoint"]