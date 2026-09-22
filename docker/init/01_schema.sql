CREATE TABLE scores (
    id    INTEGER PRIMARY KEY,
    value INTEGER NOT NULL
);

INSERT INTO scores (id, value) VALUES
    (1, 120),
    (2, 80),
    (3, -10);

CREATE TABLE boxes (
    id     INTEGER PRIMARY KEY,
    number INTEGER NOT NULL
);

INSERT INTO boxes (id, number) VALUES
    (1, 5),
    (2, 2),
    (3, 4);
