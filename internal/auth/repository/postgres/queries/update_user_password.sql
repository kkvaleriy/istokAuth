UPDATE users
SET "passHash" = @passHash 
WHERE "UUID" = @UUID;