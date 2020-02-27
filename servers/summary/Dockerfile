FROM alpine
COPY summary /summary
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
EXPOSE 5001
ENTRYPOINT [ "/summary" ]