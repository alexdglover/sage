package models

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const Submitted string = "SUBMITTED"
const Processing string = "PROCESSING"
const Failed string = "FAILED"
const Completed string = "COMPLETED"

type ImportSubmission struct {
	gorm.Model
	FileName             string
	SubmissionDateTime   string
	Status               string
	AccountType          string
	TransactionsImported int
	TransactionsSkipped  int
	BalancesImported     int
	BalancesSkipped      int
	AccountID            uint
	Account              Account
}

type ImportSubmissionRepository struct {
	DB *gorm.DB
}

// Save is an UPSERT operation, returning the ID of the record and an optional error
func (isr *ImportSubmissionRepository) Save(submission ImportSubmission) (id uint, err error) {
	result := isr.DB.Save(&submission).Clauses(clause.Returning{Columns: []clause.Column{{Name: "id"}}})

	return submission.ID, result.Error
}
