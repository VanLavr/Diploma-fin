-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

create table if not exists groups(
    id serial primary key,
    name text
);

create table if not exists students(
    uuid text primary key,
    first_name text,
    last_name text,
    middle_name text,
    group_id integer references groups(id)
);

create table if not exists teachers(
    uuid text primary key,
    first_name text,
    last_name text,
    middle_name text,
    email text
);

create table if not exists exams(
    id serial primary key,
    name text
);

create table if not exists debts(
    id serial primary key,
    exam_id integer references exams(id),
    student_uuid text references students(uuid),
    teacher_uuid text references teachers(uuid),
    date timestamp
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table debts;
drop table teachers;
drop table students;
drop table exams;
drop table groups;
-- +goose StatementEnd
