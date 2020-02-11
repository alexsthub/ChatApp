docker rm -f gateway

docker pull alexsthub/gateway

docker run \
  -d \
  -e ADDR=:443 \
  -p 443:443 \
  -v /etc/letsencrypt:/etc/letsencrypt:ro \
  -e TLSKEY="/etc/letsencrypt/live/api.alexst.me/privkey.pem" \
  -e TLSCERT="/etc/letsencrypt/live/api.alexst.me/fullchain.pem" \
  -e SESSIONKEY="mytestkey123" \
  -e REDISADDR="127.0.0.1:6379" \
  -e DSN="root:password@tcp(127.0.0.1:3306)/users" \ 
  --name gateway \
  alexsthub/gateway

exit