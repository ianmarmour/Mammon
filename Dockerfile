FROM golang:alpine AS builder
LABEL maintainer="Ian Armour <ianmarmour@gmail.com>"

RUN apk update && apk add --no-cache git openssh make bash
RUN export GIT_TERMINAL_PROMPT=1

WORKDIR /go/src/github.com/ianmarmour/Mammon
COPY . .

RUN make

FROM alpine:latest
RUN apk update && apk add --no-cache git openssh make bash

RUN addgroup -S mammon && adduser -s /bin/bash -S -G mammon mammon 

RUN mkdir -p /home/mammon/cache
RUN mkdir -p /home/mammon/db
RUN mkdir -p /home/mammon/scripts

COPY --from=builder /go/src/github.com/ianmarmour/Mammon/bin/mammon /usr/local/bin/mammon
COPY --from=builder /go/src/github.com/ianmarmour/Mammon/scripts/entry_point.sh /home/mammon/scripts/entry_point.sh

RUN chmod +x /home/mammon/scripts/entry_point.sh
RUN chown -R mammon /home/mammon

USER mammon
WORKDIR /home/mammon


ENTRYPOINT ["/home/mammon/scripts/entry_point.sh"]