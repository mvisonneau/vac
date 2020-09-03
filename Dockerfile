##
# BUILD CONTAINER
##

FROM goreleaser/goreleaser:v0.142.0 as builder

WORKDIR /build

COPY . .
RUN \
apk add --no-cache make ca-certificates ;\
make build-linux-amd64

##
# RELEASE CONTAINER
##

FROM busybox:1.32.0-glibc

WORKDIR /

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /build/dist/vac_linux_amd64/vac /usr/local/bin/

# Run as nobody user
USER 65534

EXPOSE 8080

ENTRYPOINT ["/usr/local/bin/vac"]
CMD [""]
