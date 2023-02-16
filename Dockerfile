FROM golang:1.20.0-alpine3.17 as builder

WORKDIR /opt/dynamicresource

COPY . .

RUN go build -o dynamic-resource


FROM alpine:3.17

ENV PORT=80

WORKDIR /opt/dynamicresource

RUN apk add dumb-init

COPY --from=builder --chmod=0755 /opt/dynamicresource/dynamic-resource /usr/local/bin/

EXPOSE 80

ENTRYPOINT ["dumb-init"]
CMD ["dynamic-resource", "./"]
