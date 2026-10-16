# The build stage
FROM golang:1.26.0-trixie AS builder
ENV GOPROXY=direct

WORKDIR /app
COPY api /app/api
COPY ops /app/ops

WORKDIR /app/api
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go build -o its-log

# The run stage
FROM golang:1.26.0-trixie AS prod
RUN apt update ; apt install -y jq

WORKDIR /app
COPY --from=builder /app/api/its-log .
COPY --from=builder /app/ops/containers/itslog/scripts/internal/entrypoint.bash .
RUN chmod 755 ./its-log
CMD ["./entrypoint.bash"]
