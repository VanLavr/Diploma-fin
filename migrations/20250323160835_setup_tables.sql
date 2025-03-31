-- +goose Up
-- +goose StatementBegin
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
    teacher_uuid text references teachers(uuid)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table groups;
drop table exams;
drop table students;
drop table teachers;
drop table debts;
-- +goose StatementEnd
