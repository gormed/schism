FROM golang:1.16-alpine

ENV GO111MODULE="on"
ENV CGO_ENABLED=1
ENV GOOS=linux
# ENV GOARCH=amd64

RUN apk add --no-cache \
  --repository=http://dl-cdn.alpinelinux.org/alpine/edge/testing \
  git bash curl sqlite g++ gcc make build-base i2c-tools linux-headers ca-certificates && \
  update-ca-certificates

RUN go get -ldflags "-s -w -extldflags '-static'" github.com/go-delve/delve/cmd/dlv
WORKDIR $GOPATH/src/gitlab.void-ptr.org/go/schism
ADD . .

RUN go build -gcflags='all=-N -l' -o /usr/local/bin/schism main.go

RUN mkdir /.cache && chmod 777 . /.cache /usr/local/bin/

CMD [ "dlv", "exec", "--accept-multiclient", "--listen=:2345", "--headless=true", "--api-version=2", "--log", "/usr/local/bin/schism" ]

COPY ./build/entrypoint.sh /
ENTRYPOINT ["/entrypoint.sh"]