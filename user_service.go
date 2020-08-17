package core

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/pquerna/otp/totp"
	"golang.org/x/crypto/bcrypt"
	"image"
	"time"
)

type UserStoreServiceImpl struct {
	Db                       *gorm.DB
	MaxAllowedInvalidAttempt uint
	InvalidAttemptWindow     time.Duration
	TOTPSecretLength         uint
}

func NewUserStoreService(db *gorm.DB, maxAttempt uint, window time.Duration) IUserStoreService {
	return &UserStoreServiceImpl{
		Db:                       db,
		MaxAllowedInvalidAttempt: maxAttempt,
		InvalidAttemptWindow:     window,
		TOTPSecretLength:         20,
	}
}

func (u *UserStoreServiceImpl) FindUserByUsername(username string) (*UserModel, error) {
	user := &UserModel{}
	findUserResult := u.Db.Find(user, "username = ?", username)
	if findUserResult.RecordNotFound() {
		return nil, errors.New("user not found")
	}
	if findUserResult.Error != nil {
		return nil, findUserResult.Error
	}
	return user, nil
}

func (u *UserStoreServiceImpl) FindUserByEmail(email string) (*UserModel, error) {
	user := &UserModel{}
	result := u.Db.Find(user, "email_address = ?", email)
	if result.RecordNotFound() {
		return nil, errors.New(fmt.Sprintf("user not found with email %s", email))
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}

func (u *UserStoreServiceImpl) FindAllUser(page uint, pageSize uint) ([]UserModel, uint, error) {
	var users []UserModel
	var total uint
	query := u.Db.Select([]string{"id", "username", "email_address"}).Model(&UserModel{})
	err := query.Limit(pageSize).Offset(pageSize * page).Find(&users).Error
	if err != nil {
		return nil, 0, err
	}
	query.Count(&total)
	return users, total, nil
}

func (u *UserStoreServiceImpl) ActivateUser(id uint) error {
	return u.updateStatus(id, false)
}

func (u *UserStoreServiceImpl) updateStatus(id uint, inactive bool) error {
	user := &UserModel{}
	user.ID = id
	updateResult := u.Db.Model(user).Update("inactive", inactive)
	if updateResult.Error != nil {
		return updateResult.Error
	} else if updateResult.RowsAffected != 1 {
		return errors.New("user not found")
	}
	return nil
}

func (u *UserStoreServiceImpl) DeactivateUser(id uint) error {
	return u.updateStatus(id, true)
}

func (u *UserStoreServiceImpl) ValidatePassword(id uint, password string) error {
	user := &UserModel{}
	user.ID = id
	findResult := u.Db.Find(user)
	if findResult.RecordNotFound() {
		return errors.New("user not found")
	}
	err := findResult.Error
	if err != nil {
		return err
	}
	if user.Inactive {
		return errors.New("user inactive")
	}
	cred := &UserCredentials{}
	credResult := u.Db.Find(cred, "user_id = ? and cred_type = ?", id, CredTypePassword)
	if credResult.RecordNotFound() {
		return errors.New("credential not found")
	}
	if credResult.Error != nil {
		return credResult.Error
	}
	if cred.Bocked {
		return errors.New("credential blocked")
	}
	err = bcrypt.CompareHashAndPassword([]byte(cred.Value), []byte(password))
	if err != nil {
		cred.IncrementInvalidAttempt(u.MaxAllowedInvalidAttempt, u.InvalidAttemptWindow)
		u.Db.Save(cred)
		return errors.New("password mismatch")
	}
	return nil
}

func (u *UserStoreServiceImpl) SetPassword(id uint, password string) (err error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), 13)
	if err != nil {
		return
	}
	return u.updateCredential(id, string(hashed), CredTypePassword)
}

func (u *UserStoreServiceImpl) GenerateTOTP(id uint, issuer string) (img image.Image, secret string, err error) {
	user := &UserModel{}
	user.ID = id
	err = u.Db.Find(user).Error
	if err != nil {
		return nil, "", err
	}
	opt := totp.GenerateOpts{
		Issuer:      issuer,
		AccountName: user.Username,
		SecretSize:  u.TOTPSecretLength,
	}
	key, err := totp.Generate(opt)
	if err != nil {
		return nil, "", err
	}
	img, err = key.Image(200, 200)
	if err != nil {
		return nil, "", err
	}
	secret = key.Secret()
	err = u.updateCredential(id, secret, CredTypeTOTP)
	if err != nil {
		return nil, "", err
	}
	return
}

func (u *UserStoreServiceImpl) ValidateTOTP(id uint, code string) error {
	user := &UserModel{}
	user.ID = id
	findResult := u.Db.Find(user)
	if findResult.RecordNotFound() {
		return errors.New("user not found")
	}
	err := findResult.Error
	if err != nil {
		return err
	}
	if user.Inactive {
		return errors.New("user inactive")
	}
	cred := &UserCredentials{}
	credResult := u.Db.Find(cred, "user_id = ? and cred_type = ?", id, CredTypeTOTP)
	if credResult.RecordNotFound() {
		return errors.New("credential not found")
	}
	if credResult.Error != nil {
		return credResult.Error
	}
	if cred.Bocked {
		return errors.New("credential blocked")
	}
	valid := totp.Validate(code, cred.Value)
	if !valid {
		cred.IncrementInvalidAttempt(u.MaxAllowedInvalidAttempt, u.InvalidAttemptWindow)
		u.Db.Save(cred)
		err = errors.New("totp validation failed")
		return err
	} else {
		return nil
	}
}

func (u *UserStoreServiceImpl) updateCredential(id uint, hashed string, credType uint8) error {
	user := &UserModel{}
	user.ID = id
	tx := u.Db.Begin()
	defer func() {
		tx.RollbackUnlessCommitted()
	}()
	findResult := tx.Find(user)
	if findResult.RecordNotFound() {
		return errors.New("user not found")
	}
	err := findResult.Error
	if err != nil {
		return err
	}
	var cred UserCredentials
	result := tx.Find(&cred, "user_id = ? and cred_type = ?", id, credType)
	if result.Error != nil && result.Error.Error() != "record not found" {
		return result.Error
	}
	if result.RecordNotFound() {
		cred = UserCredentials{
			UserID: id,
			Type:   credType,
			Value:  hashed,
			Bocked: false,
		}
	} else {
		cred.Value = hashed
	}
	err = tx.Save(&cred).Error
	return err
}
