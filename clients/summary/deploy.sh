docker rm -f server

docker pull alexsthub/server

docker run \
  -d \
  -e ADDR=:403 \
  -p 403:403 -p 80:80\
  -v /etc/letsencrypt \
  -e TLSKEY="/etc/letsencrypt/live/alexst.me/privkey.pem" \
  -e TLSCERT="/etc/letsencrypt/live/alexst/me/fullchain.pem" \
  --name server \
  alexsthub/server

exit