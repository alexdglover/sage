-- balances definition

DROP TABLE balances_v2;

CREATE TABLE balances_v2 (
id int primary key,
date text not null,
effective_start_date text not null,
effective_end_date text,
amount real not null,
account_id int not null,
FOREIGN KEY (account_id) REFERENCES accounts (id) ON DELETE CASCADE 
ON UPDATE NO ACTION
);

INSERT INTO balances_v2 (
	id,
	date,
	effective_start_date,
	effective_end_date,
	amount,
	account_id
) VALUES (
	1,
	datetime('now'), -- This will be the actual datetime from the statement in UTC
	date('now'),
	date('now', '+180 days'),
	544548.97,
	1
);

INSERT INTO balances_v2 (
	id,
	date,
	effective_start_date,
	effective_end_date,
	amount,
	account_id
) VALUES (
	2,
	datetime('now'), -- This will be the actual datetime from the statement in UTC
	date('now'),
	date('now', '+1 year'),
	544544.97,
	1
);

INSERT INTO balances_v2 (
	id,
	date,
	effective_start_date,
	effective_end_date,
	amount,
	account_id
) VALUES (
	3,
	datetime('now'), -- This will be the actual datetime from the statement in UTC
	date('now'),
	null,
	544.11,
	1
);
INSERT INTO balances_v2 (
	id,
	date,
	effective_start_date,
	effective_end_date,
	amount,
	account_id
) VALUES (
	4,
	datetime('now'), -- This will be the actual datetime from the statement in UTC
	date('now'),
	date('now', '+180 days'),
	544548.97,
	2
);

INSERT INTO balances_v2 (
	id,
	date,
	effective_start_date,
	effective_end_date,
	amount,
	account_id
) VALUES (
	5,
	datetime('now'), -- This will be the actual datetime from the statement in UTC
	date('now'),
	date('now', '+1 year'),
	544589.97,
	2
);

INSERT INTO balances_v2 (
	id,
	date,
	effective_start_date,
	effective_end_date,
	amount,
	account_id
) VALUES (
	6,
	datetime('now'), -- This will be the actual datetime from the statement in UTC
	date('now'),
	null,
	602.11,
	2
);
