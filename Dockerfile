FROM golang:1.13 AS build-env

# Add namespace here to resolve /vendor dependencies
ENV NAMESPACE github.com/schjan/fritzdect-exporter
WORKDIR /go/src/$NAMESPACE

ADD . ./

ARG version=dev

RUN CGO_ENABLED=0 go build -v -ldflags "-w -s"  -a -installsuffix cgo -o /exporter main.go


FROM scratch
LABEL maintainer="j.schaefer@estwx.de"

EXPOSE 8000
COPY --from=build-env /usr/local/go/lib/time/zoneinfo.zip /usr/local/go/lib/time/zoneinfo.zip
COPY --from=build-env /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=build-env /exporter /

ENTRYPOINT [ "./exporter" ]