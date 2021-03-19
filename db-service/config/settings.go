package config

// Settings ...
type Settings struct {
	DbType     string `json:"DBTYPE"`
	DbHost     string `json:"DBHOST" `
	DbDataBase string `json:"DBDATABASE" `
	DbUsername string `json:"DBUSERNAME"`
	DbPassword string `json:"DBPASSWORD"`
	DbPort     string `json:"DBPORT"`
	DbSchema   string `json:"DBSCHEMA"`
	DbSSLMode  string `json:"DBSSLMODE"`
}
