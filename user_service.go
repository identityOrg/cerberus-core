package core

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/identityOrg/cerberus-core/models"
	"github.com/identityOrg/oidcsdk"
	"github.com/jinzhu/gorm"
	"github.com/pquerna/otp/totp"
	"golang.org/x/crypto/bcrypt"
	"image"
)

type UserStoreServiceImpl struct {
	Db     *gorm.DB
	Config *Config
}

func NewUserStoreServiceImpl(db *gorm.DB, config *Config) *UserStoreServiceImpl {
	return &UserStoreServiceImpl{Db: db, Config: config}
}

func (u *UserStoreServiceImpl) FindUserByUsername(ctx context.Context, username string) (*models.UserModel, error) {
	user := &models.UserModel{}
	db := u.Db
	findUserResult := db.Find(user, "username = ?", username)
	if findUserResult.RecordNotFound() {
		return nil, errors.New("user not found")
	}
	if findUserResult.Error != nil {
		return nil, findUserResult.Error
	}
	return user, nil
}

func (u *UserStoreServiceImpl) FindUserByEmail(ctx context.Context, email string) (*models.UserModel, error) {
	user := &models.UserModel{}
	db := u.Db
	result := db.Find(user, "email_address = ?", email)
	if result.RecordNotFound() {
		return nil, errors.New(fmt.Sprintf("user not found with email %s", email))
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}

func (u *UserStoreServiceImpl) FindAllUser(ctx context.Context, page uint, pageSize uint) ([]models.UserModel, uint, error) {
	var users []models.UserModel
	var total uint
	db := u.Db
	query := db.Select([]string{"id", "username", "email_address"}).Model(&models.UserModel{})
	err := query.Limit(pageSize).Offset(pageSize * page).Find(&users).Error
	if err != nil {
		return nil, 0, err
	}
	query.Count(&total)
	return users, total, nil
}

func (u *UserStoreServiceImpl) ActivateUser(ctx context.Context, id uint) error {
	return u.updateStatus(ctx, id, false)
}

func (u *UserStoreServiceImpl) updateStatus(ctx context.Context, id uint, inactive bool) error {
	user := &models.UserModel{}
	user.ID = id
	db := u.Db
	updateResult := db.Model(user).Update("inactive", inactive)
	if updateResult.Error != nil {
		return updateResult.Error
	} else if updateResult.RowsAffected != 1 {
		return errors.New("user not found")
	}
	return nil
}

func (u *UserStoreServiceImpl) DeactivateUser(ctx context.Context, id uint) error {
	return u.updateStatus(ctx, id, true)
}

func (u *UserStoreServiceImpl) ValidatePassword(ctx context.Context, id uint, password string) error {
	user := &models.UserModel{}
	user.ID = id
	db := u.Db
	findResult := db.Find(user)
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
	cred := &models.UserCredentials{}
	credResult := db.Find(cred, "user_id = ? and cred_type = ?", id, CredTypePassword)
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
		cred.IncrementInvalidAttempt(u.Config.MaxInvalidLoginAttempt, u.Config.InvalidAttemptWindow)
		db.Save(cred)
		return errors.New("password mismatch")
	}
	return nil
}

func (u *UserStoreServiceImpl) SetPassword(ctx context.Context, id uint, password string) (err error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), u.Config.PasswordCost)
	if err != nil {
		return
	}
	return u.updateCredential(ctx, id, string(hashed), CredTypePassword)
}

func (u *UserStoreServiceImpl) GenerateTOTP(ctx context.Context, id uint, issuer string) (image.Image, string, error) {
	user := &models.UserModel{}
	user.ID = id
	db := u.Db
	result := db.Find(user)
	if result.RecordNotFound() {
		return nil, "", errors.New("user not found")
	}
	if result.Error != nil {
		return nil, "", result.Error
	}
	opt := totp.GenerateOpts{
		Issuer:      issuer,
		AccountName: user.Username,
		SecretSize:  u.Config.TOTPSecretLength,
	}
	key, err := totp.Generate(opt)
	if err != nil {
		return nil, "", err
	}
	img, err := key.Image(200, 200)
	if err != nil {
		return nil, "", err
	}
	err = u.updateCredential(ctx, id, key.Secret(), CredTypeTOTP)
	if err != nil {
		return nil, "", err
	}
	return img, key.Secret(), nil
}

func (u *UserStoreServiceImpl) ValidateTOTP(ctx context.Context, id uint, code string) error {
	user := &models.UserModel{}
	user.ID = id
	db := u.Db
	findResult := db.Find(user)
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
	cred := &models.UserCredentials{}
	credResult := db.Find(cred, "user_id = ? and cred_type = ?", id, CredTypeTOTP)
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
		cred.IncrementInvalidAttempt(u.Config.MaxInvalidLoginAttempt, u.Config.InvalidAttemptWindow)
		db.Save(cred)
		err = errors.New("totp validation failed")
		return err
	} else {
		return nil
	}
}

func (u *UserStoreServiceImpl) updateCredential(ctx context.Context, id uint, hashed string, credType uint8) error {
	user := &models.UserModel{}
	user.ID = id
	db := u.Db
	findResult := db.Find(user)
	if findResult.RecordNotFound() {
		return errors.New("user not found")
	}
	err := findResult.Error
	if err != nil {
		return err
	}
	var cred models.UserCredentials
	result := db.Find(&cred, "user_id = ? and cred_type = ?", id, credType)
	if result.Error != nil && result.Error.Error() != "record not found" {
		return result.Error
	}
	if result.RecordNotFound() {
		cred = models.UserCredentials{
			UserID: id,
			Type:   credType,
			Value:  hashed,
			Bocked: false,
		}
	} else {
		cred.Value = hashed
		cred.FirstInvalidAttempt = nil
		cred.Bocked = false
	}
	err = db.Save(&cred).Error
	return err
}

func (u *UserStoreServiceImpl) GetUser(ctx context.Context, id uint) (*models.UserModel, error) {
	user := &models.UserModel{}
	user.ID = id
	db := u.Db
	findResult := db.Find(user)
	if findResult.RecordNotFound() {
		return nil, errors.New("user not found")
	}
	if findResult.Error != nil {
		return nil, findResult.Error
	}
	return user, nil
}

func (u *UserStoreServiceImpl) UsernameAvailable(ctx context.Context, username string) (available bool) {
	user := &models.UserModel{}
	db := u.Db
	result := db.Select("username").Find(user, "username = ?", username)
	return !result.RecordNotFound()
}

func (u *UserStoreServiceImpl) ChangeUsername(ctx context.Context, id uint, username string) (err error) {
	user := &models.UserModel{}
	user.ID = id
	db := u.Db
	findResult := db.Select("username").Find(user)
	if findResult.RecordNotFound() {
		return errors.New("user not found")
	}
	if findResult.Error != nil {
		return findResult.Error
	}
	updateResult := db.Model(user).Update("username", username)
	if updateResult.Error != nil {
		return updateResult.Error
	}
	return nil
}

func (u *UserStoreServiceImpl) InitiateEmailChange(ctx context.Context, id uint, email string) (code string, err error) {
	user := &models.UserModel{}
	user.ID = id
	db := u.Db
	findResult := db.Find(user)
	if findResult.RecordNotFound() {
		return "", errors.New("user not found")
	}
	if findResult.Error != nil {
		return "", findResult.Error
	}
	user.TempEmailAddress = email
	updateResult := db.Save(user)
	if updateResult.Error != nil {
		return "", updateResult.Error
	}
	if updateResult.RowsAffected != 1 {
		return "", errors.New("update email initiation failed")
	}
	return u.GenerateUserOTP(ctx, id, 6)
}

func (u *UserStoreServiceImpl) CompleteEmailChange(ctx context.Context, id uint, code string) error {
	err := u.ValidateOTP(ctx, id, code)
	if err != nil {
		return err
	}
	user := &models.UserModel{}
	user.ID = id
	db := u.Db
	findResult := db.Find(user)
	if findResult.RecordNotFound() {
		return errors.New("user not found")
	}
	if findResult.Error != nil {
		return findResult.Error
	}
	user.EmailAddress = user.TempEmailAddress
	user.TempEmailAddress = ""
	updateResult := db.Save(user)
	if updateResult.Error != nil {
		return updateResult.Error
	}
	if updateResult.RowsAffected != 1 {
		return errors.New("update email complete failed")
	}
	return nil
}

func (u *UserStoreServiceImpl) CreateUser(ctx context.Context, username string, email string, metadata *models.UserMetadata) (id uint, err error) {
	user := &models.UserModel{
		Username:         username,
		TempEmailAddress: email,
		Metadata:         metadata,
		Inactive:         true,
	}
	db := u.Db
	saveResult := db.Save(user)
	if saveResult.Error != nil {
		return 0, saveResult.Error
	}
	return user.ID, nil
}

func (u *UserStoreServiceImpl) UpdateUser(ctx context.Context, id uint, metadata *models.UserMetadata) (err error) {
	user := &models.UserModel{}
	user.ID = id
	db := u.Db
	findResult := db.Find(user)
	if findResult.RecordNotFound() {
		return errors.New("user not found")
	}
	if findResult.Error != nil {
		return findResult.Error
	}
	user.Metadata = metadata
	return db.Save(user).Error
}

func (u *UserStoreServiceImpl) PatchUser(ctx context.Context, id uint, metadata *models.UserMetadata) (err error) {
	user := &models.UserModel{}
	user.ID = id
	db := u.Db
	findResult := db.Find(user)
	if findResult.RecordNotFound() {
		return errors.New("user not found")
	}
	if findResult.Error != nil {
		return findResult.Error
	}
	jsonData, err := json.Marshal(metadata)
	if err != nil {
		return err
	}
	err = json.Unmarshal(jsonData, user.Metadata)
	if err != nil {
		return err
	}
	return db.Save(user).Error
}

func (u *UserStoreServiceImpl) DeleteUser(ctx context.Context, id uint) (err error) {
	user := &models.UserModel{}
	user.ID = id
	db := u.Db
	return db.Delete(user).Error
}

func (u *UserStoreServiceImpl) GenerateUserOTP(ctx context.Context, id uint, length uint8) (code string, err error) {
	random, err := GenerateRandom(true, length)
	if err != nil {
		return "", err
	}
	otp := &models.UserOTP{
		ValueHash: random,
		UserID:    id,
	}
	db := u.Db
	err = db.Save(otp).Error
	if err != nil {
		return "", err
	}
	return random, nil
}

func (u *UserStoreServiceImpl) ValidateOTP(ctx context.Context, id uint, code string) (err error) {
	otp := &models.UserOTP{}
	db := u.Db
	findResult := db.Find(otp, "user_id = ? and hash_value = ?", id, code)
	if findResult.Error != nil {
		return findResult.Error
	}
	return db.Delete(otp).Error
}

func (u *UserStoreServiceImpl) Authenticate(ctx context.Context, username string, credential []byte) (err error) {
	user, err := u.FindUserByUsername(ctx, username)
	if err != nil {
		return err
	}
	return u.ValidatePassword(ctx, user.ID, string(credential))
}

func (u *UserStoreServiceImpl) GetClaims(ctx context.Context, username string, scopes oidcsdk.Arguments, claimsIDs []string) (map[string]interface{}, error) {
	db := u.Db

	openid := &models.ScopeModel{
		Name:        "openid",
		Description: "Open id scope",
	}
	userNClaim := &models.ClaimModel{
		Name:        "username",
		Description: "User Name",
	}
	db.Save(openid).Association("Claims").Append(userNClaim)

	user := &models.UserModel{}
	findUser := db.Find(user, "username = ?", username)
	if findUser.Error != nil {
		return nil, findUser.Error
	}
	m := *user.Metadata

	responseMap := make(map[string]interface{})
	for _, scope := range scopes {
		sc := &models.ScopeModel{}
		result := db.Find(sc, "name = ?", scope)
		if result.RecordNotFound() {
			continue
		} else if result.Error != nil {
			return nil, result.Error
		}
		var claims []models.ClaimModel
		assResult := result.Association("Claims").Find(&claims)
		if assResult.Error != nil {
			continue
		}
		for _, cl := range claims {
			if v, ok := m[cl.Name]; ok {
				responseMap[cl.Name] = v
			}
		}
	}
	for _, ci := range claimsIDs {
		if v, ok := m[ci]; ok {
			responseMap[ci] = v
		}
	}

	return responseMap, nil
}

func (u *UserStoreServiceImpl) IsConsentRequired(context.Context, string, string, oidcsdk.Arguments) bool {
	return false
}

func (u *UserStoreServiceImpl) StoreConsent(context.Context, string, string, oidcsdk.Arguments) error {
	return nil
}

func (u *UserStoreServiceImpl) FetchUserProfile(_ context.Context, username string) oidcsdk.RequestProfile {
	profile := oidcsdk.RequestProfile{}
	profile.SetUsername(username)
	profile.SetDomain("demo.com")
	return profile
}
