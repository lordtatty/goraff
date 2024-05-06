package llm

func boolptr(b bool) *bool {
	return &b
}

type ChatOpts struct {
	StreamCh chan string
}
