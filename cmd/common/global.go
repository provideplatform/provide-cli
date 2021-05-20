package common

import (
	"github.com/provideservices/provide-go/api/ident"
	"github.com/provideservices/provide-go/common"
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

func EtherscanBaseURL(networkID string) *string {
	switch networkID {
	case "deca2436-21ba-4ff5-b225-ad1b0b2f5c59":
		return common.StringOrNil("https://etherscan.io")
	case "07102258-5e49-480e-86af-6d0c3260827d":
		return common.StringOrNil("https://rinkeby.etherscan.io")
	case "66d44f30-9092-4182-a3c4-bc02736d6ae5":
		return common.StringOrNil("https://ropsten.etherscan.io")
	case "8d31bf48-df6b-4a71-9d7c-3cb291111e27":
		return common.StringOrNil("https://kovan.etherscan.io")
	case "1b16996e-3595-4985-816c-043345d22f8c":
		return common.StringOrNil("https://goerli.etherscan.io")
	default:
		return nil
	}
}
