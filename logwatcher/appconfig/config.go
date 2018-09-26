package appconfig

const (
	LogInterval int    = 60 // Minutes
	LogPath     string = []string{"/home/ubuntu/log", "/var/log/productionlog"}
	LogFile     string = "production.log"
	ESDomain    string = "https://search-something-es-xxxxxxxxxxxxxxx.us-east-1.es.amazonaws.com" // ElasticSearch
	ESIndex     string = "Production_log"
	ESType      string = "fatal"
)
