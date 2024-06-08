package conectorBD

import "database/sql"

func conexionDB() (conexion *sql.DB) {
	Driver := "mysql"
	User := "root"
	Password := "root"
	Name := "sistemabd"

	conexion, err := sql.Open(Driver, User+":"+Password+"@tcp(127.0.0.1)/"+Name)
	if err != nil {
		panic(err.Error())
	}
	return conexion
}
