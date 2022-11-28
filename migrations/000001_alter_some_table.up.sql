create table if not exists users(
    id serial primary key,
    first_name varchar(255) not null,
    last_name varchar(255) not null,
    phone_number varchar(100) unique,
    email varchar(100) not null unique,
    image_url varchar(255),
    password varchar(100),
    created_at timestamp with time zone default current_timestamp,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone
);

create table if not exists notes(
    id serial primary key,
    user_id integer not null REFERENCES users(id),
    title varchar(255) not null,
    description text not null,
    created_at timestamp with time zone default current_timestamp,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone
);