FROM alpine:3.21.3
WORKDIR /
COPY ./build /usr/local/bin/updategateway
RUN chmod +x /usr/local/bin/updategateway
ENTRYPOINT ["/usr/local/bin/updategateway"]