package common

import (
	"github.com/provideservices/provide-go/api/ident"
)

var (
	Application   *ident.Application
	ApplicationID string

	OrganizationID string
	Organization   *ident.Organization

	UserID string
	User   *ident.User

	AccountID   string
	ConnectorID string
	ContractID  string
	NetworkID   string
	NodeID      string
	WalletID    string

	Verbose bool
)
