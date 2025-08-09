SELECT  "UUID", nickname, "userType", "isActive"
	FROM public.users
	WHERE (phone = @phone or email = @email) and passHash = @passHash;