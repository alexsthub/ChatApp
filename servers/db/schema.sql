create table if not exists users
(
  assignedID int not null auto_increment primary key,
  id int not null,
  email varchar(40) not null,
  first_name varchar(64) not null,
  last_name varchar(128) not null,
  pass_hash char(64) not null,
  username varchar(255) not null,
  photo_url varchar(255) not null,
  UNIQUE KEY unique_email (email)
  UNIQUE KEY unique_username (username)
);
