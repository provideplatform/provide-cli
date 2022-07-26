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
	"time"

	uuid "github.com/kthomas/go.uuid"

	"github.com/provideplatform/provide-go/api/baseline"
	"github.com/provideplatform/provide-go/api/ident"
)

// Workgroup is a baseline workgroup context; called WorkgroupType because Workgroup is already declared in common
type WorkgroupType struct {
	baseline.Workgroup
	Config *WorkgroupConfig `json:"config"`
}

// WorkgroupConfig is a baseline workgroup configuration object
type WorkgroupConfig struct {
	Environment        *string      `json:"environment"`
	L2NetworkID        *uuid.UUID   `json:"l2_network_id"`
	OnboardingComplete bool         `json:"onboarding_complete"`
	SystemSecretIDs    []*uuid.UUID `json:"system_secret_ids"`
	VaultID            *uuid.UUID   `json:"vault_id"`
	WebhookSecret      *string      `json:"webhook_secret"`
}

// Organization model; called OrganizationType because Organization is already declared in common
type OrganizationType struct {
	ident.Organization
	Metadata *OrganizationMetadata `json:"metadata"`
}

// Organization metadata
type OrganizationMetadata struct {
	Address    string                                       `json:"address"`
	Workgroups map[uuid.UUID]*OrganizationWorkgroupMetadata `json:"workgroups"`
}

// Organization workgroup metadata
type OrganizationWorkgroupMetadata struct {
	OperatorSeparationDegree uint32                  `json:"operator_separation_degree"`
	Privacy                  *WorkgroupMetadataLegal `json:"privacy,omitempty"`
	SystemSecretIDs          []*uuid.UUID            `json:"system_secret_ids"`
	TOS                      *WorkgroupMetadataLegal `json:"tos,omitempty"`
	VaultID                  *uuid.UUID              `json:"vault_id"`
}

// Organization workgroup metadata legal data
type WorkgroupMetadataLegal struct {
	AgreedAt  *time.Time `json:"agreed_at"`
	Signature *string    `json:"signature"`
}
