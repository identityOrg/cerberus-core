package cerberus_models

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"os"
	"testing"
)

func TestServiceProviderModel_Migrate(t *testing.T) {
	db, err := gorm.Open("sqlite3", "sp.db")
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
