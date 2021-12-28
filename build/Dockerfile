FROM golang:1.15-alpine3.14 as builder

RUN mkdir /build
WORKDIR $GOPATH/src/gitlab.void-ptr.org/go/schism
ADD . .

RUN set -ex \
  && go build \
		-gcflags "all=-trimpath=${GOPATH}" \
		-asmflags "all=-trimpath=${GOPATH}" \
    -o /build/schism main.go

FROM alpine:3.15

COPY --from=builder /build/schism /usr/local/bin/schism
RUN groupadd -g 2000 go \
  && useradd -m -u 2001 -g go go
USER go

ENTRYPOINT [ "/usr/local/bin/schism" ]