create table if not exists products (
  id serial primary key,
  name varchar(255) not null,
  description text not null,
  image varchar(255) not null,
  price decimal(10, 2) not null,
  quantity int not null,
  created_at timestamp not null default current_timestamp
);