/*
 * Copyright 2017-2022 Provide Technologies Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package common

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	provide "github.com/provideplatform/provide-go/api"
	"github.com/provideplatform/provide-go/api/ident"
	"github.com/provideplatform/provide-go/common"
)

const releaseRepositoryPackageName = "Provide"
const releaseRepositorySSHURL = "git@github.com:provideplatform/provide.git"
const releaseRepositoryHTTPSURL = "https://github.com/provideplatform/provide"
const DefaultPage = 1
const DefaultRpp = 25

var (
	ApplicationID string
	Application   *ident.Application

	WorkgroupID string
	Workgroup   *WorkgroupType

	OrganizationID string
	Organization   *OrganizationType

	UserID string
	User   *ident.User

	SubjectAccountID string

	AccountID   string
	ConnectorID string
	ContractID  string
	NetworkID   string
	L2NetworkID string
	NodeID      string
	WalletID    string

	Manifest *provide.Manifest
	Verbose  bool

	Email    string
	Password string
)

func init() {
	resolveReleaseContext()
}

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

// resolveReleaseContext attempts to parse a Provide release manifest.json
func resolveReleaseContext() {
	path := fmt.Sprintf("./manifest.json")
	if _, err := os.Stat(path); err == nil {
		manifestJSON, err := os.ReadFile(path)
		if err != nil {
			return
		}
		err = json.Unmarshal(manifestJSON, &Manifest)
	}
}

// IsReleaseContext returns true if `prvd` is run when `pwd` is the root of a Provide release
func IsReleaseContext() bool {
	if Manifest != nil {
		return Manifest.Name == releaseRepositoryPackageName && (strings.ToLower(Manifest.Repository) == releaseRepositoryHTTPSURL || strings.ToLower(Manifest.Repository) == releaseRepositorySSHURL)
	}

	return false
}

// IsReleaseRepositoryContext is not yet used...
func IsReleaseRepositoryContext() bool {
	path := fmt.Sprintf("./.git/config")
	if _, err := os.Stat(path); err == nil {
		cfg, err := os.ReadFile(path)
		if err != nil {
			return false
		}
		cfgstr := string(cfg)
		return strings.Contains(cfgstr, releaseRepositorySSHURL) || strings.Contains(cfgstr, releaseRepositoryHTTPSURL)
	}

	return false
}

// SHA256 is a convenience method to return the sha256 hash of the given input
func SHA256(str string) string {
	digest := sha256.New()
	digest.Write([]byte(str))
	return hex.EncodeToString(digest.Sum(nil))
}
