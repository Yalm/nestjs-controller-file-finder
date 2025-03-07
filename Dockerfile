FROM alpine:3.21.3
COPY ./build /bin/updateapigateway
RUN chmod +x updateapigateway
ENTRYPOINT ["/bin/updateapigateway"]