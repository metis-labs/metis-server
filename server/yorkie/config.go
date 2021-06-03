package yorkie

// Config is the configuration for Yorkie.
type Config struct {
	RPCAddr      string `json:"RPCAddr"`
	WebhookToken string `json:"WebhookToken"`
	Collection   string `json:"Collection"`
}
