FROM golang:1.15.3 AS build
WORKDIR /src
ADD go.mod go.sum ./
RUN mkdir -p ./cmd/kube-sherlock
ADD ./cmd/kube-sherlock/main.go ./config/config.yaml ./cmd/kube-sherlock
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-w' ./cmd/kube-sherlock

FROM alpine:3.7
COPY --from=build src/cmd/kube-sherlock/config.yaml app/config.yaml
COPY --from=build src/kube-sherlock app/kube-sherlock
WORKDIR /app
ENTRYPOINT ["./kube-sherlock", "labels", "-t"]

# Metadata
ARG BUILD_DATE
ARG VCS_REF
LABEL org.label-schema.build-date=$BUILD_DATE \
    org.label-schema.name="kube-sherlock" \
    org.label-schema.description="Check if labels are applied to your containers" \
    org.label-schema.url="https://github.com/cmendible/kube-sherlock" \
    org.label-schema.vcs-ref=$VCS_REF \
    org.label-schema.vcs-url="https://github.com/cmendible/kube-sherlock" \
    org.label-schema.schema-version="0.1"