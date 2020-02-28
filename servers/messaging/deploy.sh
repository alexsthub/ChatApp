docker rm -f messages
docker pull alexsthub/messages

# docker run \
#   -d \
#   --network api \
#   -p 27017:27017 \
#   --name mongoMessages \
#   mongo

docker run \
  -d \
  --network api \
  --name messages \
  -e ADDR=:6000 \
  alexsthub/messages:latest
