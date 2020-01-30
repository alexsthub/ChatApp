GOOS=linux go build
docker build -t alexsthub/server .
go clean

docker push alexsthub/server

ssh -i ~/.ssh/info441.pem ec2-user@ec2-44-231-29-154.us-west-2.compute.amazonaws.com < deploy.sh