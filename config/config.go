package config

var APP_VERSION = "dev"

var (
	MISTRAL_TOKEN     = Env("MISTRAL_TOKEN")
	BEARER_TOKEN      = Env("BEARER_TOKEN")
	RATE_LIMIT        = EnvAsInt("RATE_LIMIT", 5)
	OPENROUTER_TOKEN  = Env("OPENROUTER_TOKEN")
	OPENAI_TOKEN      = Env("OPENAI_TOKEN")
	CHUTES_API_TOKEN  = Env("CHUTES_API_TOKEN")
	GITHUB_TOKEN      = Env("GITHUB_TOKEN")
	POSTGRES_HOST     = Env("POSTGRES_HOST")
	POSTGRES_PORT     = Env("POSTGRES_PORT")
	POSTGRES_USER     = Env("POSTGRES_USER")
	POSTGRES_PASSWORD = Env("POSTGRES_PASSWORD")
	POSTGRES_DB       = Env("POSTGRES_DB")
	SERVER_PORT       = EnvWithDefault("SERVER_PORT", "9111")
	SQLITE_DB_PATH    = EnvWithDefault("SQLITE_DB_PATH", "data/content-alchemist.db")

	// README_MIN_CONTENT_LENGTH is the minimum number of "meaningful" characters
	// (after stripping links, images, HTML and markdown markup) a README must have
	// before it is sent to the LLM. Repositories below this threshold are skipped
	// to avoid generating and storing garbage/refusal descriptions.
	README_MIN_CONTENT_LENGTH = EnvAsInt("README_MIN_CONTENT_LENGTH", 150)
)
