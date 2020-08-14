package cerberus_models

import (
	"errors"
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

func NewUserStoreServiceImpl(db *gorm.DB, maxAttempt uint, window time.Duration) IUserStoreService {
	return &UserStoreServiceImpl{
		Db:                       db,
		MaxAllowedInvalidAttempt: maxAttempt,
		InvalidAttemptWindow:     window,
		TOTPSecretLength:         20,
	}
}

func (u *UserStoreServiceImpl) FindUserByUsername(username string) (*UserModel, error) {
	user := &UserModel{}
	err := u.Db.Find(user, "username = ?", username).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UserStoreServiceImpl) FindUserByEmail(email string) (*UserModel, error) {
	user := &UserModel{}
	err := u.Db.Find(user, "email_address = ?", email).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UserStoreServiceImpl) FindAllUser(page uint, pageSize uint) ([]UserModel, uint, error) {
	var users []UserModel
	var total uint
	query := u.Db.Select([]string{"id", "username", "email_address"})
	err := query.Limit(pageSize).Offset(pageSize * page).Find(&users).Error
	if err != nil {
		return nil, 0, err
	}
	query.Count(&total)
	return users, total, nil
}

func (u *UserStoreServiceImpl) ActivateUser(id uint) error {
	user := &UserModel{
		BaseModel: BaseModel{
			ID: id,
		},
	}
	return u.Db.Model(user).Update("inactive", false).Error
}

func (u *UserStoreServiceImpl) DeactivateUser(id uint) error {
	user := &UserModel{
		BaseModel: BaseModel{
			ID: id,
		},
	}
	return u.Db.Model(user).Update("inactive", true).Error
}

func (u *UserStoreServiceImpl) ValidatePassword(id uint, password string) (valid bool, err error) {
	valid = false
	err = nil
	user := &UserModel{
		BaseModel: BaseModel{
			ID: id,
		},
	}
	err = u.Db.Preload("Credentials").Find(user).Error
	if err != nil {
		return
	}
	if user.Inactive {
		err = errors.New("user inactive")
		return
	}
	for _, credential := range user.Credentials {
		if credential.Type != CredTypePassword {
			continue
		}
		if credential.Bocked {
			err = errors.New("user inactive")
			return
		}
		err = bcrypt.CompareHashAndPassword([]byte(credential.Value), []byte(password))
		if err != nil {
			credential.IncrementInvalidAttempt(u.MaxAllowedInvalidAttempt, u.InvalidAttemptWindow)
			u.Db.Save(user)
		} else {
			valid = true
		}
		return
	}
	err = errors.New("no suitable credentials found")
	return
}

func (u *UserStoreServiceImpl) SetPassword(id uint, password string) (err error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), 13)
	if err != nil {
		return
	}
	return u.updateCredential(id, string(hashed), CredTypePassword)
}

func (u *UserStoreServiceImpl) GenerateTOTP(id uint, issuer string) (img image.Image, secret string, err error) {
	user := &UserModel{
		BaseModel: BaseModel{
			ID: id,
		},
	}
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

func (u *UserStoreServiceImpl) ValidateTOTP(id uint, code string) (valid bool, err error) {
	valid = false
	err = nil
	user := &UserModel{
		BaseModel: BaseModel{
			ID: id,
		},
	}
	err = u.Db.Preload("Credentials").Find(user).Error
	if err != nil {
		return
	}
	if user.Inactive {
		err = errors.New("user inactive")
		return
	}
	for _, credential := range user.Credentials {
		if credential.Type != CredTypeTOTP {
			continue
		}
		if credential.Bocked {
			err = errors.New("credential is blocked")
			return
		}
		valid = totp.Validate(credential.Value, code)
		if !valid {
			err = errors.New("totp validation failed")
			credential.IncrementInvalidAttempt(u.MaxAllowedInvalidAttempt, u.InvalidAttemptWindow)
			u.Db.Save(user)
		}
		return
	}
	err = errors.New("no suitable credentials found")
	return
}

func (u *UserStoreServiceImpl) updateCredential(id uint, hashed string, credType uint8) (err error) {
	user := &UserModel{
		BaseModel: BaseModel{
			ID: id,
		},
	}
	tx := u.Db.Begin()
	defer func() {
		tx.RollbackUnlessCommitted()
	}()
	err = tx.Preload("Credentials").Find(user).Error
	if err != nil {
		return
	}
	if user.Inactive {
		err = errors.New("user inactive")
		return
	}
	updated := false
	for _, cred := range user.Credentials {
		if cred.Bocked {
			cred.Bocked = false
		}
		if cred.Type != credType {
			continue
		}
		cred.Value = hashed
		updated = true
	}
	if updated {
		err = tx.Save(user).Error
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	} else {
		newCred := UserCredentials{
			Type:  credType,
			Value: hashed,
		}
		err = tx.Model(user).Association("Credentials").Append(newCred).Error
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}
	return
}
