-- +goose Up
-- +goose StatementBegin
-- CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Insert teacher
-- INSERT INTO teachers(
    -- uuid,
    -- first_name,
    -- last_name,
    -- middle_name,
    -- email
-- ) VALUES 
    -- ('admin', 'ADMIN', 'ADMIN', 'ADMIN', 'admin');

-- Insert exam
-- INSERT INTO exams(
    -- id,
    -- name
-- ) VALUES 
    -- (1, 'fizicheskaia kultura');

-- Insert group
-- INSERT INTO groups(
    -- name
-- ) VALUES 
    -- ('bsbo-01-21');

-- Insert student
-- INSERT INTO students(
    -- uuid,
    -- first_name,
    -- last_name,
    -- middle_name,
    -- email,
    -- group_id
-- ) VALUES 
    -- (uuid_generate_v4(), 'Ivan', 'Lavrushko', 'Evgenievich', 'ahaha@mail.com', 1);

-- Insert debt (FIXED: removed trailing comma)
-- INSERT INTO debts(
    -- exam_id,
    -- student_uuid,
    -- teacher_uuid
-- ) VALUES 
    -- (1, (SELECT uuid FROM students WHERE first_name = 'Ivan'), (SELECT uuid FROM teachers WHERE first_name = 'Ivan'));
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- DELETE FROM debts;
-- DELETE FROM students;
-- DELETE FROM groups;
-- DELETE FROM exams;
DELETE FROM teachers;
-- +goose StatementEnd