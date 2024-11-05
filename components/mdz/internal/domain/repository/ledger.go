package repository

import (
	"github.com/LerianStudio/midaz/common/mmodel"
)

type Ledger interface {
	Create(organizationID string, inp mmodel.CreateLedgerInput) (*mmodel.Ledger, error)
	Get(organizationID string, limit, page int) (*mmodel.Ledgers, error)
	GetByID(organizationID, ledgerID string) (*mmodel.Ledger, error)
	Update(organizationID, ledgerID string, inp mmodel.UpdateLedgerInput) (*mmodel.Ledger, error)
}
