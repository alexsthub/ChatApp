docker rm -f summary

docker pull alexsthub/summary

docker run \
  -d \
  -e ADDR=:443 \
  -p 443:443 -p 80:80\
  -v /etc/letsencrypt:/etc/letsencrypt:ro \
  -e TLSKEY="/etc/letsencrypt/live/alexst.me/privkey.pem" \
  -e TLSCERT="/etc/letsencrypt/live/alexst/me/fullchain.pem" \
  --name summary \
  alexsthub/summary

exit