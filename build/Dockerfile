FROM golang:1.16-alpine3.14 as builder

ENV CGO_ENABLED=1
ENV GOOS=linux
# ENV GOARCH=amd64

RUN apk update && apk add git g++ gcc make build-base i2c-tools linux-headers

RUN mkdir /build
WORKDIR $GOPATH/src/gitlab.void-ptr.org/go/schism
ADD . .
RUN set -ex \
  && go build \
  -gcflags "all=-trimpath=${GOPATH}" \
  -asmflags "all=-trimpath=${GOPATH}" \
  -o /build/schism main.go

FROM alpine:3.15

RUN apk update && apk add sqlite

COPY --from=builder /build/schism /usr/local/bin/schism

ENV \
  USER_GO=go \
  GROUP_GO=go

RUN addgroup -S $GROUP_GO \
  && adduser -S $USER_GO -G $GROUP_GO

RUN mkdir -p /db && touch /db/schism.sqlite && chown -R $USER_GO:$GROUP_GO /db

USER $USER_GO

ENTRYPOINT [ "/usr/local/bin/schism" ]

# Setup the healthcheck
# COPY ./healthcheck.sh /healthcheck.sh
# HEALTHCHECK --interval=10s --start-period=30s --retries=3 \
#         CMD ["/healthcheck.sh"]