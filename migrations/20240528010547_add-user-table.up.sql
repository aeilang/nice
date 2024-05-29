create table if not exists users (
  id serial primary key,
  firstname varchar(255) not null,
  lastname varchar(255) not null,
  email varchar(255) unique not null,
  password varchar(255) not null,
  created_at timestamp not null default current_timestamp
);