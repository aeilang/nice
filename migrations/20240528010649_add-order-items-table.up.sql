create table if not exists order_items (
  id serial primary key,
  order_id int not null references orders(id),
  product_id int not null references products(id),
  quantity int not null,
  price decimal(10, 2) not null
);