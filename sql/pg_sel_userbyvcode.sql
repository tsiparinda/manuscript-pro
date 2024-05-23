select id, "name", email, password_hash, is_verified, verification_code
from users where verification_code = ($1);