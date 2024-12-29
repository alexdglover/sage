package models

import "gorm.io/gorm"

type Balance struct {
	gorm.Model
	EffectiveDate      string
	Amount             int
	AccountID          uint
	Account            Account
	ImportSubmissionID *uint
	ImportSubmission   *ImportSubmission
}
