FROM golang:1.23 AS builder
COPY ./monitor /etc/monitor
WORKDIR /etc/monitor
RUN go mod tidy
RUN go build -o /bin/monitor

COPY --from=coredns/coredns:1.12.0 /coredns /bin/coredns

COPY entrypoint.sh /entrypoint.sh
COPY Corefile /etc/hallucinet/Corefile
RUN chmod +x /entrypoint.sh
ENTRYPOINT ["/entrypoint.sh"]
