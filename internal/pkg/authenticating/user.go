package authenticating

import (
	"context"

	"golang.org/x/crypto/bcrypt"
)

type (
	contextKey string

	// User represents the model of a user
	User struct {
		// gorm.Model
		UID        int    //`gorm:"unique"`
		UserName   string `json:"username" gorm:"unique"`
		Password   string `json:"password"`
		Gender     string `json:"gender"`
		FirstName  string `json:"firstname"`
		LastName   string `json:"lastname"`
		MiddleName string `json:"middlename"`
		// Status           string        `json:"status"`
		// Role             string        `json:"role"`
		// Languages        []Language    `gorm:"many2many:user_languages"`
		// BillingAddress   Address       `json:"billingAddress"` // One-To-One relationship (belongs to - use BillingAddressID as foreign key)
		// BillingAddressID sql.NullInt64 `json:"billingAddressID"`
		// MailingAddress   Address       `json:"mailingAddress"` // One-To-One relationship (belongs to - use MailingAddressID as foreign key)
		// MailingAddressID sql.NullInt64 `json:"mailingAddressID"`
		// Contact          []Contact     `json:"contact"`
		// // Profile          Profile       `json:"profile"`
		// ProfileID      uint
		// EmailAddresses []EmailAddress `json:"emailAddress" gorm:"unique"`
		// DateofBirth    time.Time      `json:"dateofBirth"`
		// DateJoined     time.Time      `json:"dateJoined"`
	}

	// EmailAddress represents email addresses
	EmailAddress struct {
		UID                   uint   `gorm:"primary_key; AUTO_INCREMENT"`
		LocalPart             string `json:"localPart" gorm:"type:varchar(100);unique_index"`
		CaseInsensitiveDomain string `json:"caseInsensitiveDomain"`
		UserID                uint   `gorm:"index"`
	}

	// Contact is the contact details of a phone number
	Contact struct {
		UID         uint   `gorm:"primary_key; AUTO_INCREMENT"`
		LineNumber  string `json:"linenumber"`
		CountryCode string `json:"countrycode"`
		AreaCode    string `json:"areacode"`
		Premfix     string `json:"prefix"`
		UserID      uint
	}

	// Address is the address of a user
	Address struct {
		ID                 uint   `gorm:"primary_key; AUTO_INCREMENT"`
		StreetAddressLine1 string `json:"streetAddress1" gorm:"type:varchar(100)"`
		StreetAddressLine2 string `json:"streetAddress2" gorm:"type:varchar(100)"`
		Country            Country
		CountryID          uint `json:"counrtyID"`
		State              State
		StateID            uint `json:"stateID"`
		City               City
		CityID             uint
		PostalCode         string `json:"postalCode"`
		Province           string `json:"province"`
	}

	// Country struct
	Country struct {
		ID   uint   `gorm:"primary_key; AUTO_INCREMENT"`
		Name string `json:"name"`
	}

	// City struct
	City struct {
		UID     uint   `gorm:"primary_key; AUTO_INCREMENT"`
		StateID uint   `json:"stateID"`
		Name    string `json:"name"`
	}

	// State struct
	State struct {
		UID       uint   `gorm:"primary_key; AUTO_INCREMENT"`
		CountryID uint   `json:"countryID"`
		Name      string `json:"name"`
	}

	// Language struct
	Language struct {
		ID   int
		Name string `gorm:"index:idx_name_code"` // Create index with name, and will create combined index if find other fields defined same name
		Code string `gorm:"index:idx_name_code"` // `unique_index` also works
	}
)

var userContextKey contextKey = "user"

// SetPassword sets user's password
func (u *User) SetPassword(p string) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}

	u.Password = string(hashedPassword)
}

// GetPassword gets the password
func (u *User) GetPassword() []byte {
	return nil
}

// IsCredentialsVerified matches given password with user's password
func (u *User) IsCredentialsVerified(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))

	return err == nil
}

// NewContext creates a new context
func (u *User) NewContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, userContextKey, u)
}

// UserFromContext gets user from context
func UserFromContext(ctx context.Context) (*User, bool) {
	u, ok := ctx.Value(userContextKey).(*User)
	return u, ok
}

// UserMustFromContext gets user from context. if can't make panic
func UserMustFromContext(ctx context.Context) *User {
	u, ok := ctx.Value(userContextKey).(*User)
	if !ok {
		panic("user can't get from request's context")
	}
	return u
}
