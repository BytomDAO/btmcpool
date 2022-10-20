package types

// Client is the general client interface
type Provider interface {
	Get() ([]string, error)
	Ensure()
	Name() string
}
