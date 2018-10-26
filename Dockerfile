FROM alpine:3.7

RUN apk add --update ca-certificates tzdata && \
    rm -rf /var/cache/apk/* /tmp/*
RUN update-ca-certificates

RUN apk --no-cache add tzdata

ADD go-app /

CMD ["/go-app","-config=/mnt/config/config.yaml"]
