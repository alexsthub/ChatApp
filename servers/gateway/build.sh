GOOS=linux go build
docker build -t alexsthub/gateway .
go clean

docker push alexsthub/gateway

ssh -i ~/.ssh/info441.pem ec2-user@ec2-35-165-73-200.us-west-2.compute.amazonaws.com < deploy.sh