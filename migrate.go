package core

import (
	"bufio"
	"context"
	"fmt"
	"github.com/identityOrg/cerberus-core/models"
	"github.com/identityOrg/oidcsdk"
	"github.com/jinzhu/gorm"
	"gopkg.in/square/go-jose.v2"
	"os"
	"strings"
)

func SetupDemoData(ormDB *gorm.DB, config *Config, sdkConfig *oidcsdk.Config, redirectUri string) error {
	fmt.Println("Creating demo client with client_id=client and client_secret=client")
	demoSp := &models.ServiceProviderModel{
		Name:         "Demo Client",
		Description:  "Demo Client",
		ClientID:     "client",
		ClientSecret: "client",
		Active:       true,
		Public:       false,
		Metadata: &models.ServiceProviderMetadata{
			RedirectUris:             []string{sdkConfig.Issuer + "/redirect"},
			Scopes:                   strings.Split("openid|offline|offline_access", "|"),
			GrantTypes:               strings.Split("authorization_code|password|refresh_token|client_credentials|implicit", "|"),
			ApplicationType:          "web",
			IdTokenSignedResponseAlg: string(jose.RS256),
		},
	}
	if redirectUri != "" {
		demoSp.Metadata.RedirectUris = append(demoSp.Metadata.RedirectUris, redirectUri)
	}
	err := ormDB.Create(demoSp).Error
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
	uid, err := userService.CreateUser(ctx, "user", "user@demo.com", metadata)
	if err != nil {
		return err
	}
	err = userService.SetPassword(ctx, uid, "user")
	if err != nil {
		return err
	}
	err = userService.ActivateUser(ctx, uid)
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
			err := ormDB.DropTableIfExists(table).Error
			if err != nil {
				return fmt.Errorf("error dropping table %s:%v", table.TableName(), err)
			}
		}
	}

	fmt.Println("creating all tables")
	for _, table := range tables {
		err := ormDB.AutoMigrate(table).Error
		if err != nil {
			return fmt.Errorf("error creating table %s:%v", table.TableName(), err)
		}
	}
	ormDB.AddUniqueIndex("idx_channel_name")
	ormDB.AddUniqueIndex("idx_alg_use")
	return InitializeDefaultScope(ormDB)
}

type dbTable interface {
	TableName() string
}
