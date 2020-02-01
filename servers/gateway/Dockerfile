# add the necessary instructions
# to create a Docker container image
# for your Go API gateway server

FROM alpine
COPY gateway /gateway
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
EXPOSE 443
ENTRYPOINT [ "/gateway" ]