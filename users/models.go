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
	Name     string `json:"name" form:"name" xml:"name" validate:"required"`
	Password string `json:"password" form:"password" xml:"password" validate:"required"`
	Age      int64  `json:"age" form:"age" xml:"age"`
}

type UserOutput struct {
	ID      int64     `json:"id"`
	Name    string    `json:"name"`
	Age     int64     `json:"age"`
	RegDate time.Time `json:"reg_date"`
}

func newOneOut(name string, password string, age int64, regDate time.Time) (uOut *UserOutput, err error) {
	u, err := newOne(name, password, age, regDate)
	if err != nil {
		return nil, err
	}

	return toOut(u), nil
}

func toOut(u *User) *UserOutput {
	return &UserOutput{
		ID:      u.ID,
		Name:    u.Name,
		Age:     u.Age,
		RegDate: u.RegDate,
	}
}

func newOne(name string, password string, age int64, regDate time.Time) (u *User, err error) {
	tmpAge := sql.NullInt64{}
	if age > 0 {
		tmpAge.Valid = true
		tmpAge.Int64 = age
	}

	conn := db.Conn()
	st, err := conn.Prepare("INSERT	INTO users(name, password, age, reg_date) VALUES(?, ?, ?, ?)")
	if err != nil {
		return nil, err
	}
	defer st.Close()

	result, err := st.Exec(name, password, tmpAge, regDate)
	if err != nil {
		return nil, err
	}

	u = new(User)
	u.ID, _ = result.LastInsertId()
	u.Name = name
	u.Password = password
	if age > 0 {
		u.Age = age
	}
	u.RegDate = regDate

	return u, nil
}

func getOneOutByID(id int64) (uOut *UserOutput, err error) {
	u, err := getOneByID(id)
	if err != nil {
		return nil, err
	}

	return toOut(u), nil
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
