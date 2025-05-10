package tools

// NATSServerTools contains all NATS server-related tool definitions
type NATSServerTools struct {
	natsURL       string
	natsCredsPath string
}

// NewNATSServerTools creates a new instance of NATSServerTools
func NewNATSServerTools(natsURL, natsCredsPath string) *NATSServerTools {
	return &NATSServerTools{
		natsURL:       natsURL,
		natsCredsPath: natsCredsPath,
	}
}
