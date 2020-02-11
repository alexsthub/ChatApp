docker rm -f mysql

docker pull alexsthub/mysql

docker run -d \
--network api \
-p 3306:3306 \
--name mysql \
-e MYSQL_ROOT_PASSWORD="databasepassword" \
-e MYSQL_DATABASE="usersDB" \
alexsthub/mysql:latest