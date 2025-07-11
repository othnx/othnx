package entities

type Config struct {
	App struct {
		// Name of the application
		Name string `json:"name"`
		// URL of the application (main domain)
		URL string `json:"url"`
		// Allowed origins for the application (CORS)
		CorsAllowedOrigins []string `json:"cors_allowed_origins"`
		// New users need to be approved by an admin
		EarlyAdoptersOnly bool `json:"early_adopters_only"`
		// Redirect URL after successful login
		RedirectAfterSuccess string `json:"redirect_after_success"`
		// Redirect URL after login error
		RedirectAfterError string `json:"redirect_after_error"`
		// API keys for the application (internal requests)
		InternalAPIKeys []string `json:"internal_api_keys"`
		// Port on which the service must run
		Port int `json:"port"`
	} `json:"app"`

	DB struct {
		// DB connection string
		PostgresURL string `json:"postgres_url"`
	} `json:"db"`

	JWT JWTConfig `json:"jwt"`

	Auth struct {
		// Providers configuration
		Providers struct {
			GitHub struct {
				Enabled      bool   `json:"enabled"`
				AppName      string `json:"app_name"`
				ClientID     string `json:"client_id"`
				ClientSecret string `json:"client_secret"`
			} `json:"github"`
			Discord struct {
				Enabled      bool   `json:"enabled"`
				AppName      string `json:"app_name"`
				ClientID     string `json:"client_id"`
				ClientSecret string `json:"client_secret"`
			} `json:"discord"`
		} `json:"providers"`
	} `json:"auth"`

	Cookies struct {
		Domain   string `json:"domain"`
		Secure   bool   `json:"secure"`
		HTTPOnly bool   `json:"http_only"`
		// SameSite cookie attribute: 1 = default, 2 = lax, 3 = strict, 4 = none
		SameSite int    `json:"same_site"`
		Path     string `json:"path"`
	} `json:"cookie"`

	User struct {
		// Roles for a user, mandatory roles are: "user" and "platform_admin"
		Roles []string `json:"roles"`
	} `json:"user"`
}

type JWTConfig struct {
	// Secret key for the JWT
	Secret string `json:"secret"`
	// Access token expiration time in minutes
	AccessTokenExpirationMin int `json:"access_token_expiration_minutes"`
	// Refresh token expiration time in days
	RefreshTokenExpirationDays int `json:"refresh_token_expiration_days"`
}
