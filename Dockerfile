FROM golang:1.7.5
WORKDIR /go/src/github.com/strongjz/leveledup-api/
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o lvl-api .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
VOLUME $GOPATH/src/github.com/strongjz/leveledup-api/config
COPY --from=0 /go/src/github.com/strongjz/leveledup-api/lvl-api .

CMD ["./lvl-api"]