SELECT  "UUID", nickname, "userType", "isActive"
	FROM users
	WHERE (phone = @phone or email = @email) and "passHash" = @passHash;