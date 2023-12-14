FROM alpine:3.18.0
RUN apk update && apk add --upgrade libcrypto3 libssl3 curl jq git

COPY homein /usr/local/bin/homein

ENTRYPOINT ["homein"]