package cerberus_models

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"os"
	"testing"
	"time"
)

var encoder *json.Encoder

func init() {
	encoder = json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
}

func TestUserCredentials_Migrate(t *testing.T) {
	db, err := gorm.Open("sqlite3", "sp.db")
	if err != nil {
		t.Error(err)
	}
	db = db.Debug()
	credential := UserCredentials{}
	model := &UserModel{
		Username:     uuid.New().String(),
		EmailAddress: uuid.New().String(),
		Credentials:  []UserCredentials{credential},
	}
	println(db.AutoMigrate(model, &credential).Error)

	db.Save(model)

	model1 := &UserModel{}
	db.Preload("Credentials").Find(model1, model.ID)

	println(encoder.Encode(model1))

	var creds []UserCredentials
	db.Model(model1).Association("Credentials").Find(&creds)
	model1.Credentials = creds

	println(encoder.Encode(model1))
}

func TestUserCredentials_IncrementInvalidAttempt(t *testing.T) {
	type fields struct {
		FirstInvalidAttempt *time.Time
		InvalidAttemptCount uint
		Bocked              bool
	}
	type args struct {
		maxAllowed uint
		window     time.Duration
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		wantBlocked bool
	}{
		{
			name: "not blocked",
			fields: fields{
				FirstInvalidAttempt: dateP(time.Now().Add(-1 * time.Minute)),
				InvalidAttemptCount: 2,
				Bocked:              false,
			},
			args: args{
				maxAllowed: 5,
				window:     2 * time.Minute,
			},
			wantBlocked: false,
		},
		{
			name: "blocked",
			fields: fields{
				FirstInvalidAttempt: dateP(time.Now().Add(-1 * time.Minute)),
				InvalidAttemptCount: 2,
				Bocked:              false,
			},
			args: args{
				maxAllowed: 2,
				window:     2 * time.Minute,
			},
			wantBlocked: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := &UserCredentials{
				FirstInvalidAttempt: tt.fields.FirstInvalidAttempt,
				InvalidAttemptCount: tt.fields.InvalidAttemptCount,
				Bocked:              tt.fields.Bocked,
			}
			if gotBlocked := uc.IncrementInvalidAttempt(tt.args.maxAllowed, tt.args.window); gotBlocked != tt.wantBlocked {
				t.Errorf("IncrementInvalidAttempt() = %v, want %v", gotBlocked, tt.wantBlocked)
			}
		})
	}
}

func dateP(add time.Time) *time.Time {
	return &add
}
