package models

import (
	"encoding/json"
	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
	"testing"
)

func TestServiceProviderModel_Migrate(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("sp.db"), &gorm.Config{})
	if err != nil {
		t.Error(err)
	}
	model := &ServiceProviderModel{
		Metadata: &ServiceProviderMetadata{},
	}
	db.AutoMigrate(model)

	model.ClientID = uuid.New().String()
	model.ClientSecret = uuid.New().String()
	model.Metadata.RedirectUris = []string{"kkkk"}

	println(db.Save(model).Error)

	println(model.ID)

	modal1 := &ServiceProviderModel{}

	db.Find(modal1, model.ID)

	println(json.NewEncoder(os.Stdout).Encode(modal1))
}

func TestPatch(t *testing.T) {
	patchData := "{\"Metadata\":{\"redirect_uris\":[\"kkkk\"]}}"
	model := &ServiceProviderModel{
		Name: "some name",
	}

	println(json.Unmarshal([]byte(patchData), model))

	println(json.NewEncoder(os.Stdout).Encode(model))
}
