#FROM golang:1.15.5
FROM scratch

COPY /bin /
EXPOSE $SERVERPORT

CMD ["/go-shortener"]