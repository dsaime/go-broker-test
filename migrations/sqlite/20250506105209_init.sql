-- +goose Up
-- +goose StatementBegin
CREATE TABLE trades_q(
    id TEXT PRIMARY KEY,
    account	TEXT NOT NULL,
    symbol	TEXT NOT NULL,
    volume	FLOAT NOT NULL,
    open	FLOAT NOT NULL,
    close	FLOAT NOT NULL,
    side	INT NOT NULL,
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE trades_q;
-- +goose StatementEnd
