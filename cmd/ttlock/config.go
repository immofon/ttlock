package main

type Config struct {
	ClientID     string `toml:"client_id"`
	ClientSecret string `toml:"client_secret"`
	Username     string `toml:"username"`
	Password     string `toml:"password"`
}
