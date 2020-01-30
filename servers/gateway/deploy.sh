docker rm -f gateway

docker pull alexsthub/gateway

docker run \
  -d \
  -e ADDR=:443 \
  -p 443:443 \
  -v /etc/letsencrypt:/etc/letsencrypt:ro \
  -e TLSKEY="/etc/letsencrypt/live/api.alexst.me/privkey.pem" \
  -e TLSCERT="/etc/letsencrypt/live/api.alexst.me/fullchain.pem" \
  --name gateway \
  alexsthub/gateway

exit