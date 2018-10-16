package inmemory

type (
	// User represents the model of a user
	User struct {
		// gorm.Model
		UID       int //`gorm:"unique"`
		UserName  string
		Password  string
		Gender    string
		FirstName string
		LastName  string
	}

	// EmailAddress represents email addresses
	EmailAddress struct {
		UID                   uint
		LocalPart             string
		CaseInsensitiveDomain string
		UserID                uint
	}

	// Contact is the contact details of a phone number
	Contact struct {
		UID         uint
		LineNumber  string
		CountryCode string
		AreaCode    string
		Premfix     string
		UserID      uint
	}

	// Address is the address of a user
	Address struct {
		ID                 uint
		StreetAddressLine1 string
		StreetAddressLine2 string
		Country            Country
		CountryID          uint
		State              State
		StateID            uint
		City               City
		CityID             uint
		PostalCode         string
		Province           string
	}

	// Country struct
	Country struct {
		ID   uint
		Name string
	}

	// City struct
	City struct {
		UID     uint
		StateID uint
		Name    string
	}

	// State struct
	State struct {
		UID       uint
		CountryID uint
		Name      string
	}

	// Language struct
	Language struct {
		ID   int
		Name string
		Code string
	}
)
