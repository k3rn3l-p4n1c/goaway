FROM golang:1.10
WORKDIR /go/src/github.com/k3rn3l-p4n1c/goaway
COPY . .
COPY config.yml /var/goaway/
RUN make goawayd
RUN make goaway
RUN mv goaway /bin/
