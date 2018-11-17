package engine

type (
	// JWTSignParser ...
	JWTSignParser interface {
		Sign(claims map[string]interface{}, secret string) (map[string]interface{}, error)
		Parse(tokenStr string, secret string) (map[string]interface{}, error)
	}

	// SecurityFactory ...
	SecurityFactory interface {
		// NewSecurityFactory ...
		NewSecurityFactory() JWTSignParser
	}
)
