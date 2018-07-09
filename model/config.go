package model

type Config struct {
	Database struct {
		Host     string `json:"host"`
		Name     string `json:"name"`
		Password string `json:"password"`
	} `json:"database"`
	Host string `json:"host"`
	Port string `json:"port"`
}
