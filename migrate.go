package core

import (
	"bufio"
	"context"
	"fmt"
	"github.com/identityOrg/cerberus-core/models"
	"github.com/identityOrg/oidcsdk"
	"gopkg.in/square/go-jose.v2"
	"gorm.io/gorm"
	"os"
	"strings"
)

func SetupDemoData(ormDB *gorm.DB, config *Config, sdkConfig *oidcsdk.Config, redirectUri string) error {
	fmt.Println("Creating demo client with client_id=client and client_secret=client")
	spMetadata := &models.ServiceProviderMetadata{
		RedirectUris:             []string{sdkConfig.Issuer + "/redirect"},
		Scopes:                   strings.Split("openid|offline|offline_access", "|"),
		GrantTypes:               strings.Split("authorization_code|password|refresh_token|client_credentials|implicit", "|"),
		ApplicationType:          "web",
		IdTokenSignedResponseAlg: string(jose.RS256),
	}
	if redirectUri != "" {
		spMetadata.RedirectUris = append(spMetadata.RedirectUris, redirectUri)
	}
	enc := NewNoOpTextEncrypt()
	spService := NewSPStoreServiceImpl(ormDB, enc, enc)
	existingSP, err := spService.FindSPByClientId(context.Background(), "client")
	if err != nil {
		spId, err := spService.CreateSP(context.Background(), "Demo Client", "Demo Client", spMetadata)
		if err != nil {
			return err
		}
		existingSP, err = spService.GetSP(context.Background(), spId)
		if err != nil {
			return err
		}
	}
	existingSP.ClientID = "client"
	existingSP.ClientSecret = "client"
	existingSP.Public = false
	err = ormDB.Save(existingSP).Error
	if err != nil {
		return err
	}
	fmt.Println("Creating demo user with username=user and password=user")

	userService := NewUserStoreServiceImpl(ormDB, config)
	metadata := &models.UserMetadata{}
	metadata.SetName("Demo User")
	metadata.SetEmail("user@demo.com")
	metadata.SetEmailVerified(true)
	ctx := context.Background()
	var uid uint
	user, err := userService.FindUserByUsername(ctx, "user")
	if err != nil {
		uid, err = userService.CreateUser(ctx, "user", "user@demo.com", metadata)
		if err != nil {
			return err
		}
	} else {
		uid = user.ID
	}
	err = userService.UpdateUser(ctx, uid, metadata)
	if err != nil {
		return err
	}
	err = userService.SetPassword(ctx, uid, "user")
	if err != nil {
		return err
	}
	err = userService.ActivateUser(ctx, uid)
	if err != nil {
		return err
	}
	fmt.Println("Creating default secret key")
	secretStore := NewSecretStoreServiceImpl(ormDB)
	_, err = secretStore.GetChannelByAlgoUse(nil, "RS256", "sig")
	if err != nil {
		_, err = secretStore.CreateChannel(nil, "default", "RS256", "sig", 30)
	}
	return err
}

func SetupDBStructure(ormDB *gorm.DB, drop bool, force bool) error {
	if drop && !force {
		fmt.Printf("Do you want to continue (Y/n): ")
		reader := bufio.NewReader(os.Stdin)
		char, _, err := reader.ReadRune()
		if err != nil {
			return err
		} else {
			switch char {
			case 'Y':
			case 'y':
				fmt.Println("continuing the migration with drop table")
			default:
				fmt.Println("Aborting the migration")
				return nil
			}
		}
	}
	scopeT := &models.ScopeModel{}
	claimT := &models.ClaimModel{}
	channelT := &models.SecretChannelModel{}
	secretT := &models.SecretModel{}
	userT := &models.UserModel{}
	credentialsT := &models.UserCredentials{}
	spT := &models.ServiceProviderModel{}
	tokensT := &models.TokensModel{}

	tables := []dbTable{scopeT, claimT, channelT, secretT, userT, credentialsT, spT, tokensT}

	fmt.Println("dropping all tables")
	if drop {
		for _, table := range tables {
			err := ormDB.Migrator().DropTable(table)
			if err != nil {
				return fmt.Errorf("error dropping table %s:%v", table.TableName(), err)
			}
		}
	}

	fmt.Println("creating all tables")
	for _, table := range tables {
		err := ormDB.AutoMigrate(table)
		if err != nil {
			return fmt.Errorf("error creating table %s:%v", table.TableName(), err)
		}
	}
	//ormDB.Model(channelT).AddUniqueIndex("idx_channel_name")//TODO fix
	//ormDB.Model(channelT).AddUniqueIndex("idx_alg_use")
	return InitializeDefaultScope(ormDB)
}

type dbTable interface {
	TableName() string
}
