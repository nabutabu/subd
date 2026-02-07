package types

type Service struct {
	Name        string
	Sub         string
	Description string
}

type State struct {
	Services map[string]Service
}
