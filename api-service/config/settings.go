package config

// Settings ...
type Settings struct {
	Port         string `json:"PORT"`
	Host         string `json:"HOST"`
	JWTSecretKey string `json:"KEY"`
	ZipkinHost   string `json:"ZIPKINHOST"`
	ZipkinPort   string `json:"ZIPKINPORT"`
}
