use users;

CREATE TABLE users.user_info (
	user_id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
	pw_hash CHAR(60) NOT NULL,
	user_name VARCHAR(100) NOT NULL COLLATE utf8mb4_general_ci,
	email VARCHAR(100) NOT NULL COLLATE utf8mb4_general_ci,
	first_name VARCHAR(100) NOT NULL COLLATE utf8mb4_general_ci,
	last_name VARCHAR(100) NOT NULL COLLATE utf8mb4_general_ci,
	INDEX user_name_index (user_name)
) DEFAULT COLLATE utf8mb4_bin;

