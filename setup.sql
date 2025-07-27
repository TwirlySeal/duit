CREATE DATABASE duit;

CREATE TABLE projects (
    id serial PRIMARY KEY NOT NULL,
    name text NOT NULL
);

INSERT INTO projects (name) VALUES
    ('Personal'),
    ('Work');

CREATE TABLE tasks (
    title text NOT NULL,
    done boolean NOT NULL,
    project_id int REFERENCES projects NOT NULL
);

INSERT INTO tasks VALUES
    ('Find Walter', false, 1),
    ('Achieve enlightenment', true, 1),
    ('Steal the moon', false, 2);
