# Initial stage: download modules
FROM golang:1.15.5 as modules

ADD go.mod go.sum /m/
RUN cd /m && go mod download

# Intermediate stage: Build the binary
FROM golang:1.15.5 as builder

COPY --from=modules /go/pkg /go/pkg

# add a non-privileged user
RUN useradd -u 10001 myapp

RUN mkdir -p /shortener
ADD . /shortener
WORKDIR /shortener

# Build the binary with go build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build /shortener
RUN chmod +x /shortener/go-shortener

# Final stage: Run the binary
FROM scratch

ENV SERVERPORT 8080
ENV URLDB redis
ENV REDISURL redis://172.28.1.4:6379
ENV REDISTIMEOUT 10

# certificates to interact with other services
# COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# don't forget /etc/passwd from previous stage
COPY --from=builder /etc/passwd /etc/passwd
USER myapp

# and finally the binary
COPY --from=builder /shortener /shortener
EXPOSE $SERVERPORT

CMD ["/shortener/go-shortener"]