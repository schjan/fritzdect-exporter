FROM --platform=$BUILDPLATFORM golang:1.14 AS build
ARG BUILDPLATFORM

# Add namespace here to resolve /vendor dependencies
ENV NAMESPACE github.com/schjan/fritzdect-exporter
WORKDIR /go/src/$NAMESPACE

ADD . ./

ARG version=dev

RUN go mod vendor -v

ARG TARGETOS
ARG TARGETARCH
ARG TARGETVARIANT

RUN GOOS=$TARGETOS GOARCH=$TARGETARCH CGO_ENABLED=0 go build -v -ldflags "-w -s"  -a -installsuffix cgo -o /exporter main.go

FROM scratch
LABEL maintainer="j.schaefer@estwx.de"

EXPOSE 8000
COPY --from=build /usr/local/go/lib/time/zoneinfo.zip /usr/local/go/lib/time/zoneinfo.zip
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=build /exporter /

ENTRYPOINT [ "./exporter" ]