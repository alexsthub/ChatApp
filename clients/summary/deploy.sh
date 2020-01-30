docker rm -f server

docker pull alexsthub/server

docker run \
  -d \
  -e ADDR=:403 \
  -p 403:403 -p 80:80\
  -v /etc/letsencrypt \
  -e TLSKEY="/etc/letsencrypt/live/your-host-name.com/privkey.pem" \
  -e TLSCERT="/etc/letsencrypt/live/your-host-name.com/fullchain.pem" \
  --name server \
  alexsthub/server

exit