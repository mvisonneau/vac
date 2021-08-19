##
# BUILD CONTAINER
##

FROM alpine:3.14.1 as certs

RUN \
  apk add --no-cache ca-certificates

##
# RELEASE CONTAINER
##

FROM busybox:1.33.1-glibc

WORKDIR /

COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY vac /usr/local/bin/

# Run as nobody user
USER 65534

EXPOSE 8080

ENTRYPOINT ["/usr/local/bin/vac"]
CMD [""]
