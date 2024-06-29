DROP TABLE IF EXISTS transactions;

CREATE TABLE transactions (
    id int primary key,
    date text not null,
    description text,
    amount real not null,
    excluded int not null, -- this will be used as a boolean for whether or not to include transactions in aggregations, like transfers
    account_id int not null,
    category_id int not null,
    FOREIGN KEY (account_id) REFERENCES accounts (id) ON DELETE CASCADE 
    ON UPDATE NO ACTION,
    FOREIGN KEY (category_id) REFERENCES categories (id) ON DELETE CASCADE 
    ON UPDATE NO ACTION
);

INSERT INTO transactions 
    (id, date, description, amount, excluded, account_id, category_id)
VALUES 
    (1, date('now'), 'Paycheck', 1241.11, false, 1, 1);

INSERT INTO transactions 
    (id, date, description, amount, excluded, account_id, category_id)
VALUES 
    (2, date('now'), 'ATM Withdrawal', 100.00, false, 1, 1);
