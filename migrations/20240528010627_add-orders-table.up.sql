create table if not exists orders (
  id serial primary key,
  user_id int not null,
  total numeric(10, 2) not null,
  status varchar(20) not null default 'pending' check (status in ('pending', 'completed', 'cancelled')),
  address text not null,
  created_at timestamp not null default current_timestamp
);