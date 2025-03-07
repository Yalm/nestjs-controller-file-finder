FROM alpine:3.21.3
RUN apk add --no-cache git openssh ca-certificates
COPY ./build /bin/updateapigateway
ENTRYPOINT ["/bin/updateapigateway"]