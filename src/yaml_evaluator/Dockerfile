# syntax=docker/dockerfile:experimental
ARG GO_VERSION=1.19
ARG ALPINE_VERSION=3.16
FROM golang:${GO_VERSION} as base

WORKDIR /workspace
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/cache \
    go mod download

FROM golang:${GO_VERSION} as builder

WORKDIR /workspace
COPY --from=base ${GOPATH}/pkg ${GOPATH}/
COPY . /workspace
RUN --mount=type=cache,target=/root/.cache/go-build \
    go build

FROM alpine:${ALPINE_VERSION} as prod
WORKDIR /workspace
COPY --from=builder /workspace/yaml_evaluator /
ENTRYPOINT /yaml_evaluator
