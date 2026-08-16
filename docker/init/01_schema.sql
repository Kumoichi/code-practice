CREATE TABLE scores (
    id    INTEGER PRIMARY KEY,
    value INTEGER NOT NULL
);

INSERT INTO scores (id, value) VALUES
    (1, 120),
    (2, 80),
    (3, -10);

CREATE TABLE dice_rolls (
    id   INTEGER PRIMARY KEY,
    pips INTEGER NOT NULL
);

INSERT INTO dice_rolls (id, pips) VALUES
    (1, 5),
    (2, 2),
    (3, 4);
