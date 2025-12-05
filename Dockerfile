# The build stage
FROM golang:1.25.5-trixie AS builder
RUN apt-get update && apt-get install -y make
WORKDIR /app
COPY . .
RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
RUN make build

# The run stage
FROM golang:1.25.5-trixie
WORKDIR /app
COPY --from=builder /app/itslog .
# We want to mount -v ${PWD}/data:/data for SQLite writing.
COPY --from=builder /app/container-config.yaml config.yaml
CMD ["./itslog"]