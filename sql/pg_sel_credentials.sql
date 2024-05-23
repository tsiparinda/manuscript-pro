select id, "name", email, is_verified, coalesce("description", '')
from users;