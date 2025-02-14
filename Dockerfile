FROM node:18-alpine AS webui
WORKDIR /etc/webui
COPY webui/package.json webui/package-lock.json .
RUN npm ci
COPY webui ./
RUN npm run build

FROM golang:1.23 AS builder

WORKDIR /etc/monitor
 
# Handle go dependencies separately so we can cache the downloads
COPY ./monitor/go.mod ./monitor/go.sum /etc/monitor
RUN go mod download

# Then build
COPY ./monitor /etc/monitor
RUN go build -o /bin/monitor

COPY --from=coredns/coredns:1.12.0 /coredns /bin/coredns
COPY --from=webui /etc/webui/dist /etc/hallucinet/webui

COPY entrypoint.sh /entrypoint.sh
COPY Corefile /etc/hallucinet/Corefile
RUN chmod +x /entrypoint.sh
ENTRYPOINT ["/entrypoint.sh"]
