ALTER TABLE users
ADD COLUMN role smallint;

ALTER TABLE users 
ADD CONSTRAINT users_role_fkey
FOREIGN KEY (role) REFERENCES roles(role_id) ON DELETE SET NULL;
