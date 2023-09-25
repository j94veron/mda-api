package main

type Configuration struct {
	DBConfiguration DBMysql
	MDAConfig       MDA
}

type DBMysql struct {
	User   string `json:"db_user"`
	Pass   string `json:"db_pass"`
	Server string `json:"db_server"`
	Port   string `json:"db_port"`
	Name   string `json:"db_name"`
}

type MDA struct {
	Path string `json:"path"`
}
