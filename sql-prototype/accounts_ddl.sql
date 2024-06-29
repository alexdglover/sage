-- accounts definition

DROP TABLE IF EXISTS accounts;

CREATE TABLE accounts (id int PRIMARY KEY, name TEXT NOT NULL, account_type TEXT NOT NULL, charge_type TEXT NOT NULL);

INSERT INTO accounts
(id, name, account_type, charge_type)
VALUES
(1, 'schwaby', 'checking', 'asset');
