package postgres

func (p *Postgres) CreateUser(name string, email string, password string, phone string, role string, address string) (int, error) {
	stmt, err := p.Db.Prepare("INSERT INTO users (name, email, password,phone, role, address) VALUES ($1, $2, $3, $4,$5, $6) RETURNING user_id")

	if err != nil {
		return 0, err
	}

	defer stmt.Close()

	var userId int
	err = stmt.QueryRow(name, email, password, phone, role, address).Scan(&userId)

	if err != nil {
		return 0, err
	}

	return int(userId), nil
}