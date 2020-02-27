GOOS=linux go build
docker build -t alexsthub/summary .
go clean

docker push alexsthub/summary

ssh -i ~/.ssh/infoKey.pem ec2-user@ec2-44-230-107-9.us-west-2.compute.amazonaws.com < deploy.sh