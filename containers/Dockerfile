# The build stage
FROM golang:1.25.5-trixie AS builder
RUN apt-get update && apt-get install -y make
ENV GOPROXY=direct
# RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
WORKDIR /tmp
ADD https://github.com/sqlc-dev/sqlc/releases/download/v1.30.0/sqlc_1.30.0_linux_amd64.tar.gz .
RUN tar xvzf sqlc_1.30.0_linux_amd64.tar.gz
RUN mv /tmp/sqlc /usr/local/bin/sqlc
WORKDIR /app
COPY .. /app
RUN make build

# The run stage
FROM golang:1.25.5-trixie AS prod
WORKDIR /app
COPY --from=builder /app/its-log .
# We want to mount -v ${PWD}/data:/data for SQLite writing.
COPY --from=builder /app/containers/container-config.yaml config.yaml
RUN chmod 755 ./its-log
CMD ["./its-log", "serve"]