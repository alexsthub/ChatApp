docker rm -f gateway
docker rm -f redisDB

docker pull alexsthub/gateway


docker run \
  -d \
  -p 6379:6379 \
  --network api \
  --name redisDB \
  redis



docker run \
  -d \
  --network api \
  -e ADDR=:443 \
  -p 443:443 \
  -v /etc/letsencrypt:/etc/letsencrypt:ro \
  -e TLSKEY="/etc/letsencrypt/live/api.alexst.me/privkey.pem" \
  -e TLSCERT="/etc/letsencrypt/live/api.alexst.me/fullchain.pem" \
  -e SESSIONKEY="mytestkey123" \
  -e REDISADDR="127.0.0.1:6379" \
  -e DSN="root:databasepassword@tcp(mysql:3306)/usersDB" \
  --name gateway \
  alexsthub/gateway

exit