package users

import (
	"database/sql"
	"fmt"
	"time"

	"echo-demo/db"

	"github.com/alexedwards/argon2id"
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
	// Use Argon2 algorithms to generate password hashes
	// Thanks Alex Edwards
	hashPass, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return nil, err
	}

	tmpAge := sql.NullInt64{}
	if age > 0 {
		tmpAge.Valid = true
		tmpAge.Int64 = age
	}

	conn := db.Conn()
	st, err := conn.Prepare("INSERT INTO users(name, password, age, reg_date) VALUES(?, ?, ?, ?)")
	if err != nil {
		return nil, err
	}
	defer st.Close()

	result, err := st.Exec(name, hashPass, tmpAge, regDate)
	if err != nil {
		return nil, err
	}

	u = new(User)
	u.ID, _ = result.LastInsertId()
	u.Name = name
	u.Password = hashPass
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

func getAllOut(limit int64, offset int64) (uOuts []*UserOutput, err error) {
	us, err := getAll(limit, offset)
	if err != nil {
		return nil, err
	}

	uOuts = make([]*UserOutput, 0, limit)
	for _, u := range us {
		uOuts = append(uOuts, toOut(u))
	}

	return uOuts, nil
}

func getAll(limit int64, offset int64) (us []*User, err error) {
	sqlStr := "SELECT id, name, password, age, reg_date FROM users"
	sqlStr += " LIMIT " + fmt.Sprintf("%d", limit)
	if offset > 0 {
		sqlStr += " OFFSET " + fmt.Sprintf("%d", offset)
	}

	conn := db.Conn()
	st, err := conn.Prepare(sqlStr)
	if err != nil {
		return nil, err
	}
	defer st.Close()

	rows, err := st.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	us = make([]*User, 0, limit)
	for rows.Next() {
		u := new(User)
		var tmpAge sql.NullInt64
		if err := rows.Scan(&u.ID, &u.Name, &u.Password, &tmpAge,
			&u.RegDate); err != nil {
			return nil, err
		}
		if tmpAge.Valid {
			u.Age = tmpAge.Int64
		}
		us = append(us, u)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return us, nil
}

func updateOneOut(id int64, name string, password string, age int64) (uOut *UserOutput, err error) {
	u, err := updateOne(id, name, password, age)
	if err != nil {
		return nil, err
	}

	return toOut(u), nil
}

func updateOne(id int64, name string, password string, age int64) (u *User, err error) {
	hashPass, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return nil, err
	}

	conn := db.Conn()
	st, err := conn.Prepare("UPDATE users SET name = ?, password = ?, age = ? WHERE id = ?")
	if err != nil {
		return nil, err
	}
	defer st.Close()

	tmpAge := sql.NullInt64{}
	if age > 0 {
		tmpAge.Valid = true
		tmpAge.Int64 = age
	}
	//result, err := st.Exec(name, hashPass, tmpAge, id)
	_, err = st.Exec(name, hashPass, tmpAge, id)
	if err != nil {
		return nil, err
	}
	//num == 0 means id not found or data unchanged.
	/*
		if num, _ := result.RowsAffected(); num == 0 {
			return nil, sql.ErrNoRows
		}
	*/

	u, err = getOneByID(id)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func deleteOne(id int64) error {
	conn := db.Conn()
	st, err := conn.Prepare("DELETE FROM users WHERE id = ?")
	if err != nil {
		return err
	}
	defer st.Close()

	result, err := st.Exec(id)
	if err != nil {
		return err
	}
	if num, _ := result.RowsAffected(); num == 0 {
		return sql.ErrNoRows
	}

	return nil
}
