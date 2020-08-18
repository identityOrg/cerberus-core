package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"gopkg.in/square/go-jose.v2"
)

type ServiceProviderModel struct {
	BaseModel
	Name         string                   `sql:"column:name;not null" json:"name,omitempty"`
	Description  string                   `sql:"column:description;size:1024" json:"description,omitempty"`
	ClientID     string                   `sql:"column:client_id;unique_index;not null" json:"client_id,omitempty"`
	ClientSecret string                   `sql:"column:client_secret" json:"client_secret,omitempty"`
	Metadata     *ServiceProviderMetadata `sql:"column:metadata;type:lob" json:"metadata,omitempty"`
}

func (sp ServiceProviderModel) TableName() string {
	return "t_sp"
}

type ServiceProviderMetadata struct {
	RedirectUris                 []string               `json:"redirect_uris,omitempty"`
	ResponseTypes                []string               `json:"response_types,omitempty"`
	GrantTypes                   []string               `json:"grant_types,omitempty"`
	ApplicationType              string                 `json:"application_type,omitempty"`
	Contacts                     []string               `json:"contacts,omitempty"`
	LogoUri                      string                 `json:"logo_uri,omitempty"`
	ClientUri                    string                 `json:"client_uri,omitempty"`
	PolicyUri                    string                 `json:"policy_uri,omitempty"`
	TosUri                       string                 `json:"tos_uri,omitempty"`
	JwksUri                      string                 `json:"jwks_uri,omitempty"`
	Jwks                         *jose.JSONWebKeySet    `json:"jwks,omitempty"`
	SectorIdentifierUri          string                 `json:"sector_identifier_uri,omitempty"`
	SubjectType                  string                 `json:"subject_type,omitempty"`
	IdTokenSignedResponseAlg     string                 `json:"id_token_signed_response_alg,omitempty"`
	IdTokenEncryptedResponseAlg  string                 `json:"id_token_encrypted_response_alg,omitempty"`
	IdTokenEncryptedResponseEnc  string                 `json:"id_token_encrypted_response_enc,omitempty"`
	UserinfoSignedResponseAlg    string                 `json:"userinfo_signed_response_alg,omitempty"`
	UserinfoEncryptedResponseAlg string                 `json:"userinfo_encrypted_response_alg,omitempty"`
	UserinfoEncryptedResponseEnc string                 `json:"userinfo_encrypted_response_enc,omitempty"`
	RequestObjectSigningAlg      string                 `json:"request_object_signing_alg,omitempty"`
	RequestObjectEncryptionAlg   string                 `json:"request_object_encryption_alg,omitempty"`
	RequestObjectEncryptionEnc   string                 `json:"request_object_encryption_enc,omitempty"`
	TokenEndpointAuthMethod      string                 `json:"token_endpoint_auth_method,omitempty"`
	TokenEndpointAuthSigningAlg  string                 `json:"token_endpoint_auth_signing_alg,omitempty"`
	DefaultMaxAge                int                    `json:"default_max_age,omitempty"`
	RequireAuthTime              bool                   `json:"require_auth_time,omitempty"`
	DefaultAcrValues             []string               `json:"default_acr_values,omitempty"`
	InitiateLoginUri             string                 `json:"initiate_login_uri,omitempty"`
	RequestUris                  []string               `json:"request_uris,omitempty"`
	OtherAttributes              map[string]interface{} `json:"other_attributes,omitempty"`
}

func (spm *ServiceProviderMetadata) Scan(src interface{}) error {
	if bytes, ok := src.([]byte); ok {
		return json.Unmarshal(bytes, spm)
	} else {
		return errors.New("failed to unmarshal value for ServiceProviderMetadata")
	}
}

func (spm *ServiceProviderMetadata) Value() (marshal driver.Value, err error) {
	marshal, err = json.Marshal(spm)
	return
}

func (spm *ServiceProviderMetadata) GetAttribute(attrName string) interface{} {
	return spm.OtherAttributes[attrName]
}

func (spm *ServiceProviderMetadata) DeleteAttribute(attrName string) {
	delete(spm.OtherAttributes, attrName)
}

func (spm *ServiceProviderMetadata) SetAttribute(attrName string, value interface{}) {
	spm.OtherAttributes[attrName] = value
}
