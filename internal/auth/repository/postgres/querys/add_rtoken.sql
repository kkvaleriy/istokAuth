INSERT INTO rtoken(
	"UUID", 
    "userUUID", 
    nickname, 
    "createdAt", 
    "expiresAt")
	VALUES (@UUID, @userUUID, @nickname, @createdAt, @expiresAt);