CREATE TABLE IF NOT EXISTS users
(
    id         TEXT PRIMARY KEY,
    username   TEXT      NOT NULL,
    created_at TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS exercises
(
    id         INTEGER PRIMARY KEY,
    code       TEXT      NOT NULL,
    name       TEXT      NOT NULL,
    created_at TIMESTAMP NOT NULL,

    UNIQUE (code)
);

CREATE INDEX IF NOT EXISTS exercises_code_idx ON exercises (code);

CREATE TABLE IF NOT EXISTS cases
(
    _id         INTEGER PRIMARY KEY AUTOINCREMENT,
    id          INTEGER   NOT NULL,
    created_at  TIMESTAMP NOT NULL,
    exercise_id INTEGER   NOT NULL REFERENCES exercises (id),
    value       BLOB      NOT NULL,

    UNIQUE (id, exercise_id)
);

CREATE INDEX IF NOT EXISTS cases_id_idx ON cases (id);

CREATE TABLE IF NOT EXISTS casualties
(
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at  TIMESTAMP NOT NULL,
    exercise_id INTEGER   NOT NULL,
    four_d      INTEGER   NOT NULL,
    case_id     INTEGER   NOT NULL REFERENCES cases (id),
    UNIQUE (exercise_id, four_d)
);

CREATE INDEX IF NOT EXISTS cadets_four_d_index ON casualties (four_d);

CREATE TABLE IF NOT EXISTS casualty_case_deterioration
(
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    casualty_id INTEGER   NOT NULL REFERENCES casualties (id) ON DELETE CASCADE,
    created_at  TIMESTAMP NOT NULL,
    value       TEXT      NOT NULL
);

CREATE INDEX IF NOT EXISTS casualties_case_deterioration_casualty_id_idx ON casualty_case_deterioration (casualty_id);

CREATE TABLE IF NOT EXISTS casualties_case_logs
(
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    casualty_id INTEGER   NOT NULL REFERENCES casualties (id) ON DELETE CASCADE,
    created_at  TIMESTAMP NOT NULL,
    type        TEXT      NOT NULL,
    data        BLOB      NOT NULL
);

CREATE INDEX IF NOT EXISTS casualty_case_logs_casualty_id_idx ON casualties_case_logs (casualty_id);

CREATE TABLE IF NOT EXISTS exercise_logs
(
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    exercise_id INTEGER   NOT NULL REFERENCES exercises (id),
    user_id     TEXT      NOT NULL REFERENCES users (id),
    type        TEXT      NOT NULL,
    created_at  TIMESTAMP NOT NULL
);

CREATE INDEX IF NOT EXISTS exercise_logs_exercise_id_idx ON exercise_logs (exercise_id);
