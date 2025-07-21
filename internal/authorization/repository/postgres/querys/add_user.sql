INSERT INTO public.users (
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
)VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);