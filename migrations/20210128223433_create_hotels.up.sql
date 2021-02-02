CREATE TABLE hotels
(
    id serial not null primary key,
    description text,
    cost int,
    creation_date date default now()
);