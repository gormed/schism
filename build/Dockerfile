FROM golang:1.15-alpine3.14 as builder

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

RUN mkdir /build
WORKDIR $GOPATH/src/gitlab.void-ptr.org/go/schism
ADD . .

RUN apk update && apk add git
RUN set -ex \
  && go build \
		-gcflags "all=-trimpath=${GOPATH}" \
		-asmflags "all=-trimpath=${GOPATH}" \
    -o /build/schism main.go

FROM alpine:3.15

COPY --from=builder /build/schism /usr/local/bin/schism

ENV \
  USER_GO=go \
  GROUP_GO=go

RUN addgroup -S $GROUP_GO \
  && adduser -S $USER_GO -G $GROUP_GO

USER $USER_GO

ENTRYPOINT [ "/usr/local/bin/schism" ]