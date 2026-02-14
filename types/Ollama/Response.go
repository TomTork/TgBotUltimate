package Ollama

type Response struct {
	Model    string `json:"model"`
	Response string `json:"response"`
	Done     bool   `json:"done"`
}
