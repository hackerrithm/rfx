package mongo

import "gopkg.in/mgo.v2/bson"

type (
	contextKey string

	// User represents the model of a user
	User struct {
		// gorm.Model
		UID       bson.ObjectId `bson:"_id" json:"id"`
		UserName  string        `bson:"username" json:"username"`
		Password  string        `bson:"password" json:"password"`
		Gender    string        `bson:"gender" json:"gender"`
		FirstName string        `bson:"firstname" json:"firstname"`
		LastName  string        `bson:"lastname" json:"lastname"`
	}

	// EmailAddress represents email addresses
	EmailAddress struct {
		UID                   bson.ObjectId `bson:"_id" json:"id"`
		LocalPart             string        `json:"localPart" gorm:"type:varchar(100);unique_index"`
		CaseInsensitiveDomain string        `json:"caseInsensitiveDomain"`
		UserID                uint          `gorm:"index"`
	}

	// Contact is the contact details of a phone number
	Contact struct {
		UID         bson.ObjectId `bson:"_id" json:"id"`
		LineNumber  string        `json:"linenumber"`
		CountryCode string        `json:"countrycode"`
		AreaCode    string        `json:"areacode"`
		Premfix     string        `json:"prefix"`
		UserID      uint
	}

	// Address is the address of a user
	Address struct {
		ID                 bson.ObjectId `bson:"_id" json:"id"`
		StreetAddressLine1 string        `json:"streetAddress1" gorm:"type:varchar(100)"`
		StreetAddressLine2 string        `json:"streetAddress2" gorm:"type:varchar(100)"`
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
		ID   bson.ObjectId `bson:"_id" json:"id"`
		Name string        `json:"name"`
	}

	// City struct
	City struct {
		UID     bson.ObjectId `bson:"_id" json:"id"`
		StateID uint          `json:"stateID"`
		Name    string        `json:"name"`
	}

	// State struct
	State struct {
		UID       bson.ObjectId `bson:"_id" json:"id"`
		CountryID uint          `json:"countryID"`
		Name      string        `json:"name"`
	}

	// Language struct
	Language struct {
		ID   int
		Name string `gorm:"index:idx_name_code"` // Create index with name, and will create combined index if find other fields defined same name
		Code string `gorm:"index:idx_name_code"` // `unique_index` also works
	}
)
