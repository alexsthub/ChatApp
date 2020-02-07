docker stop mysql
docker rm -f mysql

docker build -t alexsthub/mysql .

docker run -d \
-p 3306:3306 \
--name mysql \
-e MYSQL_ROOT_PASSWORD="databasepassword" \
alexsthub/mysql:latest
