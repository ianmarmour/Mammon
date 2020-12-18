FROM golang:alpine AS builder
LABEL maintainer="Ian Armour <ianmarmour@gmail.com>"

RUN apk update && apk add --no-cache git openssh make

RUN export GIT_TERMINAL_PROMPT=1

WORKDIR /go/src/github.com/ianmarmour/Mammon
COPY . .

RUN make

FROM scratch
COPY --from=builder /go/src/github.com/ianmarmour/Mammon/bin/mammon /go/bin/mammon
ENTRYPOINT ["/go/bin/mammon"]