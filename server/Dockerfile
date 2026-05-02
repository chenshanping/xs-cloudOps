FROM golang:1.24.12-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/server .

FROM alpine:3.22 AS runtime

RUN apk add --no-cache ca-certificates tzdata curl

WORKDIR /app

COPY --from=builder /out/server ./server
COPY config.docker.yaml ./config.docker.yaml
COPY rbac_model.conf ./rbac_model.conf

RUN mkdir -p /app/uploads /app/logs

EXPOSE 9000

CMD ["./server", "-c", "config.docker.yaml"]
