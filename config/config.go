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

	// README_MAX_NON_LATIN_PERCENT is the maximum share (in percent) of non-Latin
	// letters a README may contain and still be treated as English. READMEs written
	// in Chinese, Japanese, Korean, Cyrillic, etc. are skipped because the pipeline
	// only publishes English-language repositories.
	README_MAX_NON_LATIN_PERCENT = EnvAsInt("README_MAX_NON_LATIN_PERCENT", 20)

	// MIN_DESCRIPTION_LENGTH is the minimum length (in runes) of a single generated
	// description segment. The generation prompt asks for a very concise description
	// (max ~275 characters), so this only guards against empty or truncated output;
	// refusal messages are caught by the refusal-phrase detector instead.
	MIN_DESCRIPTION_LENGTH = EnvAsInt("MIN_DESCRIPTION_LENGTH", 40)
)
