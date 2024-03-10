package middleware

type (
	// KeyAuthConfig defines the config for KeyAuth middleware.
	KeyAuthConfig struct {

		// KeyLookup is a string in the form of "<source>:<name>" that is used
		// to extract key from the request.
		// Optional. Default value "header:Authorization".
		// Possible values:
		// - "header:<name>"
		// - "query:<name>"
		// - "form:<name>"
		KeyLookup string `header:Authorization"`

		// AuthScheme to be used in the Authorization header.
		// Optional. Default value "Bearer".
		AuthScheme string
	}

)

var (
	// DefaultKeyAuthConfig is the default KeyAuth middleware config.
	DefaultKeyAuthConfig = KeyAuthConfig{
		AuthScheme: "Bearer",
	}
)
