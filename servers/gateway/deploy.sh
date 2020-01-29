docker rm -f gateway

docker pull alexsthub/gateway

docker run \
  -d \
  -e ADDR=:403 \
  -p 403:403 \
  -v /etc/letsencrypt \
  -e TLSKEY="/etc/letsencrypt/live/your-host-name.com/privkey.pem" \
  -e TLSCERT="/etc/letsencrypt/live/your-host-name.com/fullchain.pem" \
  --name gateway \
  alexsthub/gateway

exit