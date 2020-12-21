FROM golang:alpine AS builder
LABEL maintainer="Ian Armour <ianmarmour@gmail.com>"

RUN apk update && apk add --no-cache git openssh make
RUN export GIT_TERMINAL_PROMPT=1

WORKDIR /go/src/github.com/ianmarmour/Mammon
COPY . .

RUN make
RUN useradd -ms /bin/bash mammon
USER mammon
WORKDIR /home/mammon

RUN mkdir -p /home/mammon/cache
RUN mkdir -p /home/mammon/db
RUN mkdir -P /home/mammon/scripts

FROM scratch
COPY --from=builder /go/src/github.com/ianmarmour/Mammon/bin/mammon /usr/local/bin/mammon
COPY --from=builder /go/src/github.com/ianmarmour/Mammon/scripts/entry_point.sh /home/mammon/scripts/entry_point.sh
RUN chmod +x /home/mammon/scripts/entry_point.sh

ENV MAMMON_CACHE_PATH=/home/mammon/cache
ENV MAMMON_DB_PATH=/home/mammon/db
ENV BLIZZARD_API_CLIENT_LOCALE=en_US
ENV BLIZZARD_API_CLIENT_REGION=us
ENV BLIZZARD_API_CLIENT_ID=/run/secrets/BLIZZARD_API_CLIENT_ID
ENV BLIZZARD_API_CLIENT_SECRET=/run/secrets/BLIZZARD_API_CLIENT_SECRET

ENTRYPOINT ["/home/mammon/scripts/entry_point.sh"]