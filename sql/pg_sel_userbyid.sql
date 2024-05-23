select id, "name", email, password_hash, is_verified, COALESCE(verification_code, '')
from users where id = ($1);