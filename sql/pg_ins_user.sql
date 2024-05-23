INSERT INTO users ("name", email, verification_code, password_hash)  
Values (($1), ($2), ($3), ($4)) RETURNING id;

