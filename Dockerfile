FROM alpine:3.21.3
WORKDIR /
COPY ./build /bin/updategateway
ENTRYPOINT ["/bin/updategateway"]