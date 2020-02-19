docker build -t alexsthub/summary .

docker push alexsthub/summary

ssh -i ~/.ssh/infoKey.pem ec2-user@ec2-44-233-24-50.us-west-2.compute.amazonaws.com < deploy.sh