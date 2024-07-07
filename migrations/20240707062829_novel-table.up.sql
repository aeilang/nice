create table if not exists novels (
    id serial primary key,
    title varchar(255) not null,
    keyword text not null,
    short text not null,
    content text not null,
    author_id int not null,
    created_at timestamp not null default current_timestamp,

    foreign key (author_id) references users(id)
);