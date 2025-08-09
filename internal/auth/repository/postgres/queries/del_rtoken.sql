DELETE FROM rtoken
WHERE "UUID" = $1
RETURNING "userUUID", nickname;