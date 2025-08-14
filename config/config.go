package config

import (
	"os"
	"strconv"
)

type Config struct {
	DSN             string
	JWTSecret       string
	AccessTokenTTL  int64
	RefreshTokenTTL int64
	ServerPort      string
}

func LoadConfig() *Config {
	c := &Config{}
	c.DSN = os.Getenv("DATABASE_DSN")
	c.JWTSecret = os.Getenv("JWT_SECRET")
	if v := os.Getenv("ACCESS_TOKEN_TTL"); v != "" {
		t, _ := strconv.ParseInt(v, 10, 64)
		c.AccessTokenTTL = t
	} else {
		c.AccessTokenTTL = 900
	}
	if v := os.Getenv("REFRESH_TOKEN_TTL"); v != "" {
		t, _ := strconv.ParseInt(v, 10, 64)
		c.RefreshTokenTTL = t
	} else {
		c.RefreshTokenTTL = 86400
	}
	c.ServerPort = os.Getenv("SERVER_PORT")
	if c.ServerPort == "" {
		c.ServerPort = "8080"
	}
	return c
}
