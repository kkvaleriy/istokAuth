package queries

import _ "embed"

var (
	//go:embed add_user.sql
	AddUser string
	//go:embed user_type.sql
	UserType string
	//go:embed update_user_password.sql
	UpdateUserPassword string
	//go:embed add_rtoken.sql
	AddRToken string
	//go:embed del_rtoken.sql
	DelRToken string
	//go:embed check_credentials.sql
	CheckUserByCredentials string
)
