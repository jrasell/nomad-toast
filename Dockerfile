FROM alpine:latest

LABEL maintainer James Rasell<(jamesrasell@gmail.com)> (@jrasell)
LABEL vendor "jrasell"

ENV NOMAD_TOAST_VERSION 0.0.1

WORKDIR /usr/bin/

RUN buildDeps=' \
                bash \
                wget \
        ' \
        set -x \
        && apk --no-cache add $buildDeps ca-certificates \
        && wget -O levant https://github.com/jrasell/nomad-toast/releases/download/${NOMAD_TOAST_VERSION}/nomad-toast_linux_amd64 \
        && chmod +x /usr/bin/nomad-toast \
        && apk del $buildDeps \
        && echo "Build complete."

CMD ["nomad-toast", "--help"]
