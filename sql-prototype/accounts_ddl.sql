-- accounts definition

DROP TABLE IF EXISTS accounts;

CREATE TABLE accounts (
    id int PRIMARY KEY, 
    name TEXT NOT NULL, 
    account_type TEXT NOT NULL, 
    charge_type TEXT NOT NULL,
    interest_rate DECIMAL(5,4) DEFAULT 0.0000
);

INSERT INTO accounts
(id, name, account_type, charge_type)
VALUES
(1, 'schwaby', 'checking', 'asset', 0.0000),
(2, 'credit_card', 'credit', 'liability', 0.1899),
(3, 'mortgage', 'mortgage', 'liability', 0.0650);
