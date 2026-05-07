# Stage 1: Build JupyterLite static assets
FROM python:3.11-slim AS builder

WORKDIR /build

RUN pip install --no-cache-dir \
    jupyter-server \
    jupyterlite-core \
    jupyterlite-pyodide-kernel \
    pandas \
    numpy \
    requests

COPY content /build/content

RUN jupyter lite build --contents content --output-dir /build/_output

FROM nginx:alpine

COPY --from=builder /build/_output /usr/share/nginx/html

EXPOSE 80