package config

type Database struct {
	Driver    string
	Host      string
	Port      string
	User      string
	Password  string
	Name      string
}

type Services struct {
	SportsMonks
}

type SportsMonks struct {
	BaseUri string
	ApiKey	string
}
