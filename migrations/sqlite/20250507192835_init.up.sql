CREATE TABLE trades_q
(
    id      TEXT PRIMARY KEY,
    account TEXT  NOT NULL,
    symbol  TEXT  NOT NULL,
    volume  FLOAT NOT NULL,
    open    FLOAT NOT NULL,
    close   FLOAT NOT NULL,
    side    TEXT   NOT NULL,
    worker_id TEXT NOT NULL,
    job_status TEXT NOT NULL,
    profit  FLOAT NOT NULL
);