package cmd

//DefaultConfig default config
var DefaultConfig = map[string]interface{}{
	"name":              "just-projecting",
	"port":              "8000",
	"log_level":         "DEBUG",
	"log_format":        "json",
	"tracer":            "no-op",
	"newrelic_apikey":   "",
	"jaeger_url":        "http://localhost:14268/api/traces",
	"mysql_dsn":         "root:benyamin@tcp(localhost:3306)/majoodb?parseTime=true",
	"grpc_port":         "8080",
	"grpc_gateway_port": "8090",
}
