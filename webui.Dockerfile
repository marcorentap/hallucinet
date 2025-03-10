FROM node:18-alpine AS webui
WORKDIR /etc/webui
COPY webui/package.json webui/package-lock.json ./
RUN npm ci
COPY webui ./
RUN npm run build

FROM nginx:1.27.4
COPY --from=webui /etc/webui/dist /etc/hallucinet/webui
