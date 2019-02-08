FROM golang:1.11.5 AS build
WORKDIR /src
ADD go.mod go.sum ./
RUN go get -v
ADD kube-sherlock.go config.yaml ./
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-w'

FROM alpine:3.7
COPY --from=build src/config.yaml app/config.yaml
COPY --from=build src/kube-sherlock app/kube-sherlock
WORKDIR /app
CMD ./kube-sherlock

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