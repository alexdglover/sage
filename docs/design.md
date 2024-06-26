# Architecture

## Basic concepts

Intuit's Mint offered reporting on several personal finance aspects:
* List of all transactions across all accounts and financial institutions
* Assets, liabilities, and net worth over time
* Monthly spending by category, or by other dimensions
* Monthly spending over time
* 

While Mint was extremely useful, you only need two types of data to achieve
all of the reporting listed above. Namely, you need
* Every individual transaction. This provides the data needed for a searchable list of transcations, reporting on income/spending by category or other dimensions, and spending/income over time
* Balances for each account. This provides the data needed for tracking assets, liabilities, and net worth over time

By using these 2 data sets, Sage is much simpler than a true accounting 
application. We don't need to worry about double entry accounting, or keeping
a perfect ledger to calculate the current balance for an account. We simplify
the app by letting the financial institutions do the hard part, and Sage
will simply use the calculated balance (as part of the data provided in
statements from financial institutions) for reporting purposes.

## Entities and Relationships

```mermaidjs
erDiagram
    ACCOUNT {
        string name
        string type
    }
    ACCOUNT ||--o{ TRANSACTION : contains
    ACCOUNT ||--o{ BALANCE : has
    TRANSACTION {
        int date
        string description
        float amount
        string excluded
        int account_fk
    }
    BALANCE {
        int effective_date
        float balance
        int account_fk
    }
    
```
