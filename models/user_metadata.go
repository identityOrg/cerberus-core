package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type (
	IUserClaims interface {
		GetName() string
		SetName(name string)
		GetGivenName() string
		SetGivenName(givenName string)
		GetFamilyName() string
		SetFamilyName(familyName string)
		GetMiddleName() string
		SetMiddleName(middleName string)
		GetNickname() string
		SetNickname(nickname string)
		GetPreferredUsername() string
		SetPreferredUsername(preferredUsername string)
		GetProfile() string
		SetProfile(profile string)
		GetPicture() string
		SetPicture(picture string)
		GetWebsite() string
		SetWebsite(website string)
		GetEmail() string
		SetEmail(email string)
		GetEmailVerified() bool
		SetEmailVerified(emailVerified bool)
		GetGender() string
		SetGender(gender string)
		GetBirthDate() string
		SetBirthDate(birthDate string)
		GetZoneInfo() string
		SetZoneInfo(zoneInfo string)
		GetLocale() string
		SetLocale(locale string)
		GetPhoneNumber() string
		SetPhoneNumber(phoneNumber string)
		GetPhoneNumberVerified() bool
		SetPhoneNumberVerified(phoneNumberVerified bool)
		GetAddress() IAddress
		SetAddress(address IAddress)
		GetAttribute(name string) interface{}
		SetAttribute(name string, value interface{})
	}
	IAddress interface {
		GetFormatted() string
		SetFormatted(formatted string)
		GetStreetAddress() string
		SetStreetAddress(streetAddress string)
		GetLocality() string
		SetLocality(locality string)
		GetRegion() string
		SetRegion(region string)
		GetPostalCode() string
		SetPostalCode(postalCode string)
		GetCountry() string
		SetCountry(country string)
		GetAttribute(name string) interface{}
		SetAttribute(name string, value interface{})
	}

	UserMetadata map[string]interface{}
	UserAddress  map[string]interface{}
)

func (u *UserMetadata) GormDBDataType(db *gorm.DB, _ *schema.Field) string {
	switch db.Dialector.Name() {
	case "mysql", "sqlite":
		return "JSON"
	case "postgres":
		return "JSONB"
	default:
		return "lob"
	}
}

func (u *UserMetadata) Scan(src interface{}) error {
	if bytes, ok := src.([]byte); ok {
		return json.Unmarshal(bytes, &u)
	} else {
		return fmt.Errorf("failed to unmarshal value for UserMetadata")
	}
}

func (u *UserMetadata) Value() (driver.Value, error) {
	return json.Marshal(u)
}

func (u UserAddress) GetFormatted() string {
	return extractStringAttribute(u, "formatted")
}

func (u UserAddress) SetFormatted(formatted string) {
	u["formatted"] = formatted
}

func (u UserAddress) GetStreetAddress() string {
	return extractStringAttribute(u, "formatted")
}

func (u UserAddress) SetStreetAddress(streetAddress string) {
	u["street_address"] = streetAddress
}

func (u UserAddress) GetLocality() string {
	return extractStringAttribute(u, "formatted")
}

func (u UserAddress) SetLocality(locality string) {
	u["locality"] = locality
}

func (u UserAddress) GetRegion() string {
	return extractStringAttribute(u, "formatted")
}

func (u UserAddress) SetRegion(region string) {
	u["region"] = region
}

func (u UserAddress) GetPostalCode() string {
	return extractStringAttribute(u, "formatted")
}

func (u UserAddress) SetPostalCode(postalCode string) {
	u["postal_code"] = postalCode
}

func (u UserAddress) GetCountry() string {
	return extractStringAttribute(u, "formatted")
}

func (u UserAddress) SetCountry(country string) {
	u["country"] = country
}

func (u UserAddress) GetAttribute(name string) interface{} {
	return u[name]
}

func (u UserAddress) SetAttribute(name string, value interface{}) {
	u[name] = value
}

func (u *UserMetadata) GetName() string {
	return extractStringAttribute(*u, "name")
}

func (u *UserMetadata) SetName(name string) {
	(*u)["name"] = name
}

func (u *UserMetadata) GetGivenName() string {
	return extractStringAttribute(*u, "given_name")
}

func (u *UserMetadata) SetGivenName(givenName string) {
	(*u)["given_name"] = givenName
}

func (u *UserMetadata) GetFamilyName() string {
	return extractStringAttribute(*u, "family_name")
}

func (u *UserMetadata) SetFamilyName(familyName string) {
	(*u)["family_name"] = familyName
}

func (u *UserMetadata) GetMiddleName() string {
	return extractStringAttribute(*u, "middle_name")
}

func (u *UserMetadata) SetMiddleName(middleName string) {
	(*u)["middle_name"] = middleName
}

func (u *UserMetadata) GetNickname() string {
	return extractStringAttribute(*u, "nickname")
}

func (u *UserMetadata) SetNickname(nickname string) {
	(*u)["nickname"] = nickname
}

func (u *UserMetadata) GetPreferredUsername() string {
	return extractStringAttribute(*u, "preferred_username")
}

func (u *UserMetadata) SetPreferredUsername(preferredUsername string) {
	(*u)["preferred_username"] = preferredUsername
}

func (u *UserMetadata) GetProfile() string {
	return extractStringAttribute(*u, "profile")
}

func (u *UserMetadata) SetProfile(profile string) {
	(*u)["profile"] = profile
}

func (u *UserMetadata) GetPicture() string {
	return extractStringAttribute(*u, "picture")
}

func (u *UserMetadata) SetPicture(picture string) {
	(*u)["picture"] = picture
}

func (u *UserMetadata) GetWebsite() string {
	return extractStringAttribute(*u, "website")
}

func (u *UserMetadata) SetWebsite(website string) {
	(*u)["website"] = website
}

func (u *UserMetadata) GetEmail() string {
	return extractStringAttribute(*u, "email")
}

func (u *UserMetadata) SetEmail(email string) {
	(*u)["email"] = email
}

func (u *UserMetadata) GetEmailVerified() bool {
	return extractBoolAttribute(u, "email_verified")
}

func (u *UserMetadata) SetEmailVerified(emailVerified bool) {
	(*u)["email_verified"] = emailVerified
}

func (u *UserMetadata) GetGender() string {
	return extractStringAttribute(*u, "gender")
}

func (u *UserMetadata) SetGender(gender string) {
	(*u)["gender"] = gender
}

func (u *UserMetadata) GetBirthDate() string {
	return extractStringAttribute(*u, "birth_date")
}

func (u *UserMetadata) SetBirthDate(birthDate string) {
	(*u)["birth_date"] = birthDate
}

func (u *UserMetadata) GetZoneInfo() string {
	return extractStringAttribute(*u, "zone_info")
}

func (u *UserMetadata) SetZoneInfo(zoneInfo string) {
	(*u)["zone_info"] = zoneInfo
}

func (u *UserMetadata) GetLocale() string {
	return extractStringAttribute(*u, "locale")
}

func (u *UserMetadata) SetLocale(locale string) {
	(*u)["locale"] = locale
}

func (u *UserMetadata) GetPhoneNumber() string {
	return extractStringAttribute(*u, "phone_number")
}

func (u *UserMetadata) SetPhoneNumber(phoneNumber string) {
	(*u)["phone_number"] = phoneNumber
}

func (u *UserMetadata) GetPhoneNumberVerified() bool {
	return extractBoolAttribute(u, "phone_number_verified")
}

func (u *UserMetadata) SetPhoneNumberVerified(phoneNumberVerified bool) {
	(*u)["phone_number_verified"] = phoneNumberVerified
}

func (u *UserMetadata) GetAddress() IAddress {
	if addr, ok := (*u)["address"].(IAddress); ok {
		return addr
	} else {
		return nil
	}
}

func (u *UserMetadata) SetAddress(address IAddress) {
	(*u)["address"] = address
}

func (u *UserMetadata) GetAttribute(name string) interface{} {
	return (*u)[name]
}

func (u *UserMetadata) SetAttribute(name string, value interface{}) {
	(*u)[name] = value
}

func extractStringAttribute(u map[string]interface{}, key string) string {
	if val, ok := u[key].(string); ok {
		return val
	} else {
		return ""
	}
}

func extractBoolAttribute(u *UserMetadata, key string) bool {
	if val, ok := (*u)[key].(bool); ok {
		return val
	} else {
		return false
	}
}
