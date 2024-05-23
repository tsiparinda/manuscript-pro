INSERT INTO users (id, "name", email, password_hash, is_verified)  
Values (($1), ($2), ($3), ($4), true);

