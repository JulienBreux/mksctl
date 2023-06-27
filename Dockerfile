ARG GO_VERSION=1.20.5
ARG APP=mksctl

FROM golang:${GO_VERSION}-alpine AS build

ARG VERSION=dev
ARG DATE=n/a
ARG COMMIT=n/a

WORKDIR /${APP}

COPY go.mod go.sum Makefile ./
COPY internal internal
COPY pkg pkg
COPY cmd cmd
COPY views views

RUN apk --no-cache add --update make libx11-dev git gcc libc-dev curl && make build

FROM gcr.io/distroless/static AS final

LABEL maintainer="Julien BREUX <julien.breux@gmail.com>"
USER nonroot:nonroot

COPY --from=build --chown=nonroot:nonroot /${APP}/bin/app /app

ENTRYPOINT ["/app"]
