FROM amazonlinux:2023

# https://github.com/grafana/k6/issues/5120
RUN dnf install -y golang git
RUN go install go.k6.io/k6@latest
RUN mv ~/go/bin/k6 /usr/local/bin/

COPY scripts /scripts

CMD [ "k6", "run", "10k.js" ]