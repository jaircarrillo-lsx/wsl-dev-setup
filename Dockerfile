# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /src

RUN apk add --no-cache git make

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /out/wsl-dev-setup .

# Runtime stage - distroless
FROM gcr.io/distroless/static:nonroot

COPY --from=builder /out/wsl-dev-setup /usr/local/bin/wsl-dev-setup

ENTRYPOINT ["wsl-dev-setup"]