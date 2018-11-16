package engine

type (
	// JWTSignParser ...
	JWTSignParser interface {
		Sign(claims map[string]interface{}, secret string) (string, error)
		Parse(tokenStr string, secret string) (map[string]interface{}, error)
	}

	// SecurityFactory ...
	SecurityFactory interface {
		// NewSecurityFactory ...
		NewSecurityFactory() SecurityFactory
	}
)
