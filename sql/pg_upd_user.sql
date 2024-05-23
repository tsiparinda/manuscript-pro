UPDATE users SET "name"=($2), email=($3), is_verified=($4), verification_code=($5), password_hash=($6)  where id=($1);
