# The build stage
FROM golang:1.26.0-trixie AS builder
RUN apt-get update && apt-get install -y make unzip awscli
ENV GOPROXY=direct

RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

# WORKDIR /tmp
# RUN curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip"
# RUN unzip awscliv2.zip
# RUN ./aws/install

# WORKDIR /tmp
# ADD https://github.com/sqlc-dev/sqlc/releases/download/v1.30.0/sqlc_1.30.0_linux_amd64.tar.gz .
# RUN tar xvzf sqlc_1.30.0_linux_amd64.tar.gz
# RUN mv /tmp/sqlc /usr/local/bin/sqlc
WORKDIR /app
COPY api /app/api
COPY ops /app/ops
COPY Makefile /app/Makefile
RUN make native

# The run stage
FROM golang:1.26.0-trixie AS prod
RUN apt update ; apt install -y jq awscli

WORKDIR /app
COPY --from=builder /app/api/its-log .
COPY --from=builder /app/ops/containers/itslog/scripts/internal/entrypoint.bash .
# We want to mount -v ${PWD}/data:/data for SQLite writing.
RUN chmod 755 ./its-log
CMD ["./entrypoint.bash"]