package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"github.com/identityOrg/oidcsdk"
	"gopkg.in/square/go-jose.v2"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type ServiceProviderModel struct {
	BaseModel
	Name         string                   `gorm:"column:name;not null" json:"name,omitempty"`
	Description  string                   `gorm:"column:description;size:1024" json:"description,omitempty"`
	ClientID     string                   `gorm:"column:client_id;index:uk_client_id,unique;not null" json:"client_id,omitempty"`
	ClientSecret string                   `gorm:"column:client_secret" json:"client_secret,omitempty"`
	Active       bool                     `gorm:"column:active" json:"active,omitempty"`
	Public       bool                     `gorm:"column:public" json:"public,omitempty"`
	Metadata     *ServiceProviderMetadata `gorm:"column:metadata" json:"metadata,omitempty"`
}

func (sp ServiceProviderModel) AutoMigrate(db gorm.Migrator) error {
	return db.AutoMigrate(&sp)
}

func (sp ServiceProviderModel) GetID() string {
	return sp.ClientID
}

func (sp ServiceProviderModel) GetSecret() string {
	return sp.ClientSecret
}

func (sp ServiceProviderModel) IsPublic() bool {
	return sp.Public
}

func (sp ServiceProviderModel) GetIDTokenSigningAlg() jose.SignatureAlgorithm {
	if sp.Metadata != nil {
		return jose.SignatureAlgorithm(sp.Metadata.IdTokenSignedResponseAlg)
	} else {
		return ""
	}
}

func (sp ServiceProviderModel) GetRedirectURIs() []string {
	if sp.Metadata != nil {
		return sp.Metadata.RedirectUris
	} else {
		return []string{}
	}
}

func (sp ServiceProviderModel) GetApprovedScopes() oidcsdk.Arguments {
	if sp.Metadata != nil {
		return sp.Metadata.Scopes
	} else {
		return []string{}
	}
}

func (sp ServiceProviderModel) GetApprovedGrantTypes() oidcsdk.Arguments {
	if sp.Metadata != nil {
		return sp.Metadata.GrantTypes
	} else {
		return []string{}
	}
}

func (sp ServiceProviderModel) TableName() string {
	return "t_sp"
}

type ServiceProviderMetadata struct {
	ClientName                   string                 `json:"client_name,omitempty"`
	RedirectUris                 []string               `json:"redirect_uris,omitempty"`
	ResponseTypes                []string               `json:"response_types,omitempty"`
	Scopes                       []string               `json:"scopes,omitempty"`
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

func (spm *ServiceProviderMetadata) GormDBDataType(db *gorm.DB, _ *schema.Field) string {
	switch db.Dialector.Name() {
	case "mysql", "sqlite":
		return "JSON"
	case "postgres":
		return "JSONB"
	default:
		return "lob"
	}
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
