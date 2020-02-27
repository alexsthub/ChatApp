docker rm -f summary

docker pull alexsthub/summary

docker run \
  -d \
  -e ADDR=:80
  --network api \
  --name summary \
  alexsthub/summary

exit