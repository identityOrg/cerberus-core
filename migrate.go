package core

import (
	"bufio"
	"context"
	"fmt"
	"github.com/identityOrg/cerberus-core/models"
	"github.com/jinzhu/gorm"
	"gopkg.in/square/go-jose.v2"
	"os"
	"strings"
	"time"
)

func MigrateDB(ormDB *gorm.DB, force bool, demo bool) error {
	if force {
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
	if force {
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
	if err := InitializeDefaultScope(ormDB); err != nil {
		return err
	}

	if demo {
		fmt.Println("Creating demo client with client_id=client and client_secret=client")
		demoSp := &models.ServiceProviderModel{
			Name:         "Demo Client",
			Description:  "Demo Client",
			ClientID:     "client",
			ClientSecret: "client",
			Active:       true,
			Public:       false,
			Metadata: &models.ServiceProviderMetadata{
				RedirectUris:             []string{"http://localhost:8080/redirect"},
				Scopes:                   strings.Split("openid|offline|offline_access", "|"),
				GrantTypes:               strings.Split("authorization_code|password|refresh_token|client_credentials|implicit", "|"),
				ApplicationType:          "web",
				IdTokenSignedResponseAlg: string(jose.RS256),
			},
		}
		err := ormDB.Create(demoSp).Error
		if err != nil {
			return err
		}
		fmt.Println("Creating demo user with username=user and password=user")
		config := &Config{
			EncryptionKey:          "very secret encryption key",
			MaxInvalidLoginAttempt: 3,
			InvalidAttemptWindow:   5 * time.Minute,
			TOTPSecretLength:       6,
		}

		userService := NewUserStoreServiceImpl(ormDB, config)
		metadata := &models.UserMetadata{}
		metadata.SetName("Demo User")
		metadata.SetEmail("user@demo.com")
		metadata.SetEmailVerified(true)
		_, err = userService.CreateUser(context.Background(), "user", "user@demo.com", metadata)
		if err != nil {
			return err
		}
	}
	return nil
}

type dbTable interface {
	TableName() string
}
