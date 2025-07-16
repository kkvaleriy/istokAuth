INSERT INTO users (
    name,
    lastname,
    nickname,
    email,
    userType,
    isActive,
    phone,
    UUID,
    passHash,
    createdAt
)VALUE (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)