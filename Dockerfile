FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /bin/api ./cmd/api
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /bin/batch ./cmd/batch

FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=builder /bin/api /bin/api
COPY --from=builder /bin/batch /bin/batch
COPY --from=builder /app/migrations /migrations

EXPOSE 8080

ENTRYPOINT ["/bin/api"]
