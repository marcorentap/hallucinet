FROM golang:1.23
COPY . /src

WORKDIR /src
RUN go mod tidy
RUN go build -o /bin/backend

CMD ["/bin/backend"]
