GOOS=linux go build
docker build -t alexsthub/gateway .
go clean

docker push alexsthub/gateway

ssh -i ~/.ssh/infoKey.pem ec2-user@ec2-44-230-107-9.us-west-2.compute.amazonaws.com < deploy.sh