package repositoryAuth

const (
	queryCreate = `
			INSERT INTO users (first_name,last_name,password,email,role) 
			VALUES ($1,$2,$3,$4,$5)
			RETURNING id
			`
	queryFindById = `
             SELECT id,first_name,last_name,password,email,role
             FROM users
             WHERE id = $1
             `
	queryFindByEmail = `
             SELECT id,first_name,last_name,password,email,role
             FROM users
             WHERE email = $1
             `
	queryUpdate = `
              UPDATE users
              SET first_name = $1, last_name = $2, password = $3, email = $4, role = $5 
              WHERE id = $6
              `
	queryDelete = `
 			  DELETE FROM users
 			  WHERE id = $1
 			  `
)
