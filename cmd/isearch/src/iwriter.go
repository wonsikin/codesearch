package src

// IWriter blabla
type IWriter struct {
	Content []string `json:"content,omitempty"`
}

func (i *IWriter) Write(p []byte) (n int, err error) {
	i.Content = append(i.Content, string(p))
	return len(p), nil
}
