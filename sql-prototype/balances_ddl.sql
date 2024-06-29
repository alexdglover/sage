-- balances definition

DROP TABLE IF EXISTS balances;

CREATE TABLE balances (
id int primary key,
date text not null, -- date from the statement itself
effective_start_date text not null,
effective_end_date text,
amount real not null,
account_id int not null,
FOREIGN KEY (account_id) REFERENCES accounts (id) ON DELETE CASCADE 
ON UPDATE NO ACTION
);

INSERT INTO balances (
	id,
	date,
	effective_start_date,
	effective_end_date,
	amount,
	account_id
) VALUES (
	1,
	date('2024-03-03'), -- This will be the actual datetime from the statement in UTC
	date('2024-03-03'),
	date('now', '+180 days'),
	544548.97,
	1
);

INSERT INTO balances (
	id,
	date,
	effective_start_date,
	effective_end_date,
	amount,
	account_id
) VALUES (
	2,
	date('2024-01-03'), -- This will be the actual datetime from the statement in UTC
	date('2024-01-03'),
	date('2024-04-03'),
	544544.97,
	1
);

INSERT INTO balances (
	id,
	date,
	effective_start_date,
	effective_end_date,
	amount,
	account_id
) VALUES (
	3,
	date('2023-09-23'), -- This will be the actual datetime from the statement in UTC
	date('2023-09-23'),
	null,
	544.11,
	1
);
INSERT INTO balances (
	id,
	date,
	effective_start_date,
	effective_end_date,
	amount,
	account_id
) VALUES (
	4,
	date('2024-05-13'), -- This will be the actual datetime from the statement in UTC
	date('2024-05-13'),
	date('now', '+180 days'),
	544548.97,
	2
);

INSERT INTO balances (
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

INSERT INTO balances (
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
