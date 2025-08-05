package querys

import _ "embed"

var (
	//go:embed add_user.sql
	AddUser string
	//go:embed add_rtoken.sql
	AddRToken string
	//go:embed check_credentials.sql
	CheckUserByCredentials string
)
