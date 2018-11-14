FROM golang:alpine as build
RUN apk add --no-cache ca-certificates
WORKDIR /docker
ADD . .

FROM scratch
COPY --from=build /etc/ssl/certs/ca-certificates.crt \
/etc/ssl/certs/ca-certificates.crt

COPY proffix-rest /px
CMD ["./px"]
