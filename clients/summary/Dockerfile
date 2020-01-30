# add the necessary instructions
# to create a Docker container image
# for your summary client using the
# NGINX web server image as the base

FROM nginx
ADD default.conf /etc/nginx/conf.d/default.conf
COPY summary /summary
# TODO: Copy the react folder?
COPY ./summary /var/www
EXPOSE 80 
ENTRYPOINT [ "/gateway" ]
