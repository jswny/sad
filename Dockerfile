FROM golang:1.15.3-alpine3.12

# Pinned versions
ARG OPENSSL_VERSION='1.1.1'
ARG OPENSSH_VERSION='8.3'

RUN apk add --no-cache \
  openssl=~${OPENSSL_VERSION} \
  openssh=~${OPENSSH_VERSION}

WORKDIR /go/src/app

COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

ENTRYPOINT ["sad"]
