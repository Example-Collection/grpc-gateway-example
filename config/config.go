package config

type DatabaseConfig struct {
	// Database Endpoint
	Endpoint string
	// User table
	User string
	// Region
	Region string
}

func Init() DatabaseConfig {
	return DatabaseConfig{
		Endpoint: "http://localhost:54000",
		User:     "grpc-gateway-example-user",
		Region:   "ap-northeast-2",
	}
}
