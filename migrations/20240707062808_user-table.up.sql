create type user_role as ENUM('admin', 'user', 'forbiden');

create table if not exists users (
    id serial primary key,
    name varchar(255) not null,
    email varchar(255) unique not null,
    password varchar(255) not null,
    role user_role not null,
    created_at timestamp not null default current_timestamp
);