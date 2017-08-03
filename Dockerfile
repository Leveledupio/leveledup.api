FROM golang:1.7.5
WORKDIR /go/src/github.com/strongjz/leveledup-api/
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o lvl-api .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=0 /go/src/github.com/strongjz/leveledup-api/lvl-api .

#ADD repositories /etc/apk/repositoriesa
RUN apk add --update python python-dev gfortran py-pip build-base
RUN pip install --upgrade --user awscli
RUN export PATH=~/.local/bin:$PATH
RUN chmod +x ~/.local/bin


# Install the new entry-point script
COPY secrets-entrypoint.sh /secrets-entrypoint.sh

# Overwrite the entry-point script
ENTRYPOINT [ "/secrets-entrypoint.sh"]

#CMD ["./lvl-api"]
