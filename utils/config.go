package utils

type AcmeConfig struct {
	Username      string `yaml:"username"`
	Password      string `yaml:"password"`
	Host          string `yaml:"host"`
	Port          int    `yaml:"port"`
	DBName        string `yaml:"dbname"`
	CustomerCount int    `yaml:"customer_count"`
	ProductCount  int    `yaml:"product_count"`
	OrderCount    int    `yaml:"order_count"`
	MaxClients    int    `yaml:"client_count"`
}
