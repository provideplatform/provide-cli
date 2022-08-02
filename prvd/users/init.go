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

package users

import (
	"fmt"
	"log"
	"os"

	"github.com/provideplatform/provide-cli/prvd/common"
	provide "github.com/provideplatform/provide-go/api/ident"

	"github.com/spf13/cobra"
)

// initCmd creates a new user
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Create a new user",
	Long:  `Create a new user in the configured ident instance; defaults to ident.provide.services.`,
	Run:   create,
}

var firstName string
var lastName string

func create(cmd *cobra.Command, args []string) {
	firstName = common.FreeInput("First Name", "", common.MandatoryValidation)
	lastName = common.FreeInput("Last Name", "", common.MandatoryValidation)

	if common.Email == "" {
		common.Email = common.FreeInput("Email", "", common.EmailValidation)
	}

	if common.Password == "" {
		common.Password = common.FreeInput("Password", "", common.MandatoryValidation)
	}

	resp, err := provide.CreateUser("", map[string]interface{}{
		"email":      common.Email,
		"password":   common.Password,
		"first_name": firstName,
		"last_name":  lastName,
	})
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	_, err = provide.Authenticate(common.Email, common.Password, "")
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	fmt.Printf("created user: %s", *resp.ID)
}

func init() {
	initCmd.Flags().StringVar(&common.Email, "email", "", "account email")
	initCmd.Flags().StringVar(&common.Password, "password", "", "account password")
	initCmd.Flags().StringVar(&firstName, "first-name", "", "account first name")
	initCmd.Flags().StringVar(&lastName, "last-name", "", "account last name")
}
