FROM --platform=$BUILDPLATFORM golang:1.21.5-alpine3.18@sha256:d8b99943fb0587b79658af03d4d4e8b57769b21dcf08a8401352a9f2a7228754 AS builder

RUN apk add --update --no-cache ca-certificates make git curl

ARG TARGETPLATFORM

WORKDIR /usr/local/src/benthos-openmeter

ARG GOPROXY

ENV CGO_ENABLED=0

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG VERSION

RUN go build -ldflags "-X main.version=${VERSION}" -o /usr/local/bin/benthos .

FROM alpine:3.19.0@sha256:51b67269f354137895d43f3b3d810bfacd3945438e94dc5ac55fdac340352f48

RUN apk add --update --no-cache ca-certificates tzdata bash

SHELL ["/bin/bash", "-c"]

# This is so we can reuse examples in development
WORKDIR /etc/benthos

COPY cloudevents.spec.json /etc/benthos/

COPY examples/http-server/input.yaml /etc/benthos/examples/http-server/input.yaml
COPY examples/http-server/output.yaml /etc/benthos/examples/http-server/output.yaml
COPY examples/kubernetes-pod-exec-time/config.yaml /etc/benthos/examples/kubernetes-pod-exec-time/config.yaml

COPY --from=builder /usr/local/bin/benthos /usr/local/bin/

CMD benthos
