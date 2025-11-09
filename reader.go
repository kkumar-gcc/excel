package excel

type Reader interface {
	Close() error
	Headers() []string
	Read() (map[string]string, error)
}
