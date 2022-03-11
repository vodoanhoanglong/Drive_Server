# build context at repo root: docker build -f Dockerfile ../..
FROM golang:1.16 AS builder

WORKDIR /app

COPY go.mod go.mod
COPY go.sum go.sum
COPY pkg pkg
COPY vendor vendor
COPY services/auth services/auth

ARG TAG
ARG GIT_COMMIT
WORKDIR /app/services/auth

RUN TAG=$TAG GIT_COMMIT=$GIT_COMMIT make build 

# stage 2: production image
FROM alpine
RUN apk add --no-cache ca-certificates

# Copy the binary to the production image from the builder stage.
COPY --from=builder /app/services/auth/auth-api /auth

RUN chmod +x /auth
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
USER appuser

# Run the web service on container startup.
CMD ["/auth"]
