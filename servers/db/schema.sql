CREATE DATABASE usersDB;

USE usersDB;

create table if not exists users
(
  id int not null auto_increment primary key,
  email varchar(40) not null,
  first_name varchar(64) not null,
  last_name varchar(128) not null,
  pass_hash char(64) not null,
  username varchar(255) not null,
  photo_url varchar(255) not null,
  UNIQUE KEY(email),
  unique key(username)
);

INSERT INTO users 
(email, first_name, last_name, pass_hash, username, photo_url) 
VALUES 
("alextan785@gmail.com", "Alex", "Tan", "ASDIASJNDIUSANDIUSADNASIUD", "alextan785", "myexamplephotourl.com/avatar/alexst"),
("test@gmail.com", "FirstTest", "LastTest", "MYPASSHASHBUTISITAGOODHASH", "TestingUser", "myexamplephotourl.com/avatar/test");