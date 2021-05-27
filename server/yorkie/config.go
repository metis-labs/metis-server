package yorkie

// Config is the configuration for Yorkie.
type Config struct {
	Addr         string `json:"Addr"`
	WebhookToken string `json:"WebhookToken"`
	Collection   string `json:"Collection"`
}
