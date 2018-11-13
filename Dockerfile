FROM golang:alpine as build
RUN apk add --no-cache ca-certificates
WORKDIR /build
ADD . .

FROM scratch
COPY --from=build /etc/ssl/certs/ca-certificates.crt \
/etc/ssl/certs/ca-certificates.crt

COPY bin/proffix-rest_v"$VERSION"_linux_amd64 /proffix-rest
CMD ["./proffix-rest"]
