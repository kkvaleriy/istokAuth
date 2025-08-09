INSERT INTO users (
    name,
    lastname,
    nickname,
    email,
    "userType",
    "isActive",
    phone,
    "UUID",
    "passHash",
    "createdAt"
)VALUES (@name, @lastname, @nickname, @email, @userType, @isActive, @phone, @UUID, @passHash, @createdAt);