docker rm -f summary

docker pull alexsthub/summary

docker run \
  -d \
  --network api \
  -e ADDR=:5000 \
  --name summary \
  alexsthub/summary

exit