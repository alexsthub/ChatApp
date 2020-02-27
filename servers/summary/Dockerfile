FROM alpine
COPY summary /summary
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
EXPOSE 5000
ENTRYPOINT [ "/summary" ]