FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY . /app/

WORKDIR /app/services/recommendation-service

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/service ./cmd/api

FROM gcr.io/distroless/static-debian12:nonroot AS production

COPY --from=builder /out/service /service

USER nonroot:nonroot
ENTRYPOINT ["/service"]
