FROM golang:alpine as build
RUN apk add --no-cache ca-certificates
WORKDIR /build
ADD . .

FROM scratch
COPY --from=build /etc/ssl/certs/ca-certificates.crt \
/etc/ssl/certs/ca-certificates.crt

COPY docker/proffix-rest /proffix-rest
CMD ["./proffix-rest"]
