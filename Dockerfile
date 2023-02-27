FROM alpine:latest

RUN set -x \
  && apk --no-cache add bash curl jq \
  && adduser -D -H appuser appuser \
  ;

COPY image-triggers /usr/local/bin
CMD ["/usr/local/bin/image-triggers"]
USER appuser
