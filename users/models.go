package users

import (
	"database/sql"
	"time"

	"echo-demo/db"
)

type User struct {
	ID       int64
	Name     string
	Password string
	Age      int64
	RegDate  time.Time
}

type UserInput struct {
	Name     string
	Password string
	Age      int64
}

type UserOutput struct {
	ID      int64     `json:"id"`
	Name    string    `json:"name"`
	Age     int64     `json:"age"`
	RegDate time.Time `json:"reg_date"`
}

func getOneOutByID(id int64) (uOut *UserOutput, err error) {
	u, err := getOneByID(id)
	if err != nil {
		return nil, err
	}

	uOut = &UserOutput{
		ID:      u.ID,
		Name:    u.Name,
		Age:     u.Age,
		RegDate: u.RegDate,
	}

	return uOut, nil
}

func getOneByID(id int64) (u *User, err error) {
	u = new(User)

	conn := db.Conn()
	st, err := conn.Prepare("SELECT id, name, password, age, reg_date FROM users WHERE id = ?")
	if err != nil {
		return nil, err
	}
	defer st.Close()

	var tmpAge sql.NullInt64
	if err := st.QueryRow(id).Scan(&u.ID, &u.Name, &u.Password, &tmpAge,
		&u.RegDate); err != nil {
		return nil, err
	}

	if tmpAge.Valid {
		u.Age = tmpAge.Int64
	}

	return u, nil
}
