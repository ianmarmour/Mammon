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
RUN chown -R mammon /home/mammon

COPY --from=builder /go/src/github.com/ianmarmour/Mammon/bin/mammon /usr/local/bin/mammon

USER mammon
WORKDIR /home/mammon

ENV BLIZZARD_API_CLIENT_LOCALE en_US
ENV BLIZZARD_API_CLIENT_REGION us
ENV MAMMON_CACHE_PATH /home/mammon/cache/
ENV MAMMON_DB_PATH /home/mammon/db/

ENTRYPOINT ["/usr/local/bin/mammon"]