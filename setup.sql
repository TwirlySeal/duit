CREATE DATABASE duit;
\c duit;

CREATE TABLE projects (
    id serial PRIMARY KEY NOT NULL,
    name text NOT NULL
);

INSERT INTO projects (name) VALUES
    ('Personal'),
    ('Work');

CREATE TABLE tasks (
    id serial PRIMARY KEY NOT NULL,
    title text NOT NULL,
    done boolean NOT NULL DEFAULT false,
    date DATE,
    time TIME,
    project_id int REFERENCES projects NOT NULL

    CONSTRAINT date_required_for_time CHECK (time IS NULL OR date IS NOT NULL)
);

INSERT INTO tasks (title, project_id) VALUES
    ('Find Walter', 1),
    ('Achieve enlightenment', 1),
    ('Steal the moon', 2);
