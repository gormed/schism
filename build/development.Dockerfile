FROM golang:1.16-alpine

ENV CGO_ENABLED=1
ENV GOOS=linux
# ENV GOARCH=amd64

RUN apk add --no-cache \
  --repository=http://dl-cdn.alpinelinux.org/alpine/edge/testing \
  git bash curl sqlite g++ gcc make build-base i2c-tools linux-headers
RUN go get -u github.com/cosmtrek/air

WORKDIR $GOPATH/src/gitlab.void-ptr.org/go/schism
ADD . .

RUN mkdir /.cache && chmod 777 . /.cache /usr/local/bin/

CMD [ "air", "-c", "build/air.conf" ]