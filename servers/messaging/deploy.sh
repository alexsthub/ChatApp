docker rm -f messages
docker pull alexsthub/messages

# docker run \
#   -d \
#   --network api \
#   -p 27017:27017 \
#   --name mongoMessages \
#   mongo

# docker run \
#   -d \
#   --hostname myrabbitmq \
#   --name rabbitmq \
#   -p 5672:5672 -p 15672:15672 \
#   --network api \
#   rabbitmq:3-management

docker run \
  -d \
  --network api \
  --name messages \
  -e ADDR=:6000 \
  alexsthub/messages:latest
