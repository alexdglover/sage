-- balances definition
--
--CREATE TABLE balances_v2 (
--id int primary key,
--date text not null,
--effective_start_date text not null,
--effective_end_date text,
--amount real not null,
--account_id int not null,
--FOREIGN KEY (account_id) REFERENCES accounts (id) ON DELETE CASCADE 
--ON UPDATE NO ACTION
--);



-- get items in a SCD that have already started
--select id, effective_start_date , (effective_start_date < date('now')) AS beforenow from balances_v2;


-- get items in a SCD that have not yet expired
--select id, effective_end_date, ((effective_end_date > date('now')) or (effective_end_date is null)) AS afternow from balances_v2;

-- alternatively - look into setting a separate `current` column, but maybe this is just more work/state to manage
--select id, effective_start_date, effective_end_date from balances_v2 WHERE is_current=TRUE ;

