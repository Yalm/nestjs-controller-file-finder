FROM alpine:3.21.3
COPY ./build /bin/updateapigateway
ENTRYPOINT ["/bin/updateapigateway"]