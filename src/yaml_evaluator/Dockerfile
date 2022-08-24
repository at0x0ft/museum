ARG GO_VERSION=1.19
ARG ALPINE_VERSION=3.16
FROM golang:${GO_VERSION} as base

WORKDIR /workspace
COPY go.mod go.sum ./
# TODO: add buildkit directive
RUN ls && go mod download

FROM golang:${GO_VERSION} as builder

WORKDIR /workspace
COPY --from=base ${GOPATH}/pkg ${GOPATH}/
COPY . /workspace
# TODO: add buildkit directive
RUN go build

FROM alpine:${ALPINE_VERSION} as prod
WORKDIR /workspace
COPY --from=builder /workspace/yaml_evaluator /
ENTRYPOINT /yaml_evaluator
