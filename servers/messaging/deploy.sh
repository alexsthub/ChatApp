docker rm -f messages
docker pull alexsthub/messages

docker run \
  -d \
  --network api \
  --name messaging \
  --e ADDR=:80 \
  alexsthub/messages:latest
