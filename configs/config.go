package configs

// Configuration is our applications configuration structure
type Configuration struct {
	// Port is application's default port
	Port string
	// StaticVariable defines variables that will not be updated
	StaticVariable string
	// ConnectionString is our database's connection string
	ConnectionString string
	// Address is the address we will use in our app
	Address string
}
