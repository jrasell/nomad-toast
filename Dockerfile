FROM golang:1.12.3-alpine3.9 AS builder

WORKDIR /go/src/github.com/jrasell/nomad-toast
RUN apk add -U make git
COPY GNUmakefile .
RUN make build-tools
COPY . .
RUN make build

FROM alpine:latest

LABEL maintainer James Rasell<(jamesrasell@gmail.com)> (@jrasell)
LABEL vendor "jrasell"

COPY --from=builder /go/src/github.com/jrasell/nomad-toast/bin/nomad-toast /usr/local/bin/
RUN apk add -U ca-certificates

ENTRYPOINT ["nomad-toast"]
CMD ["--help"]
