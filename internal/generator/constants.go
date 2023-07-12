package generator

const (
	Norton = "norton"
	Avast  = "avast"
	Avg    = "avg"
)

var BrandToCountries = map[string][]string{
	"avast":  {"US", "SA", "CN", "TW", "CZ", "BR", "DE", "DK", "ES", "FI", "FR", "GR", "HU", "ID", "IL", "IT", "JP", "KR", "MY", "NL", "NO", "PL", "PT", "RU", "SE", "SK", "TH", "TR", "UA"},
	"smb":    {"CZ", "DK", "NL", "US", "FR", "DE", "IT", "JP", "NO", "PL", "BR", "RU", "ES", "SE"},
	"norton": {"US", "FI", "PL", "BR", "IT", "NO", "DK", "NL", "FR", "DE", "ES", "SE"},
	"avg":    {"US", "SA", "CN", "TW", "CZ", "DK", "NL", "FI", "FR", "DE", "GR", "IL", "HU", "ID", "IT", "JP", "KR", "MY", "NO", "PL", "BR", "PT", "RU", "RS", "SK", "ES", "SE", "TH", "TR", "VN"},
}

var EN_COUNTRIES = []string{"US", "CA", "AU", "GB"}
