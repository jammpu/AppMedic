package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/sessions"
)

// variable global para la conexion a base de datos
var cDB *sql.DB

var store = sessions.NewCookieStore([]byte("secret-key"))

// llamar plantillas
var plantillas = template.Must(template.ParseGlob("src/templates/*"))

func main() {
	cDB = conexionDB()

	// Cargar archivos estaticos
	fs := http.FileServer(http.Dir("./src/styles"))
	http.Handle("/styles/", http.StripPrefix("/styles/", fs))
	// Manejadores para las rutas
	http.HandleFunc("/", StartLogin)
	http.HandleFunc("/login", login)
	http.HandleFunc("/sistema", sistema)
	http.HandleFunc("/logout", logout)

	fmt.Println("Servidor Corriendo..")
	http.ListenAndServe(":8080", nil)
}

// conexion base de datos

// plantillas
func StartLogin(w http.ResponseWriter, r *http.Request) {
	plantillas.ExecuteTemplate(w, "login", nil)
}

func login(w http.ResponseWriter, r *http.Request) {
	user := r.FormValue("user")
	password := r.FormValue("password")

	// Verificar la autenticación del usuario
	if authenticateUser(user, password) {
		// Autenticación exitosa, almacenar el nombre de usuario en la sesión
		session, _ := store.Get(r, "session-name")
		session.Values["username"] = user
		session.Save(r, w)
		// Autenticación exitosa, redirigir al usuario a la página del sistema
		http.Redirect(w, r, "/sistema", http.StatusFound)
		return
	}

	// Autenticación fallida, mostrar mensaje de error
	http.Redirect(w, r, "/", http.StatusFound)
}

// acciones
func authenticateUser(user, password string) bool {

	var storedUsername string
	var storedPassword string

	// Obtener contraseña almacenada para el usuario
	err := cDB.QueryRow("SELECT user, password FROM admin WHERE user = ?", user).Scan(&storedUsername, &storedPassword)
	if err != nil {
		log.Println(err)
		return false
	}
	if user != storedUsername {
		return false
	}

	// Comparar contraseña almacenada con la contraseña proporcionada
	err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(password))
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}

// sesion activa de usuario
func sistema(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session-name")
	username, ok := session.Values["username"].(string)
	if !ok || username == "" {
		// Usuario no autenticado, redirigir al formulario de inicio de sesión
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	// Obtener los datos del usuario de la base de datos utilizando el nombre de usuario
	userData := getUserData(username)

	// Usuario autenticado, mostrar página del sistema
	plantillas.ExecuteTemplate(w, "sistema", userData)
}

// cerrar sesion usuario
func logout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session-name")

	// Establecer la duración de la sesión en cero para cerrarla
	session.Options.MaxAge = -1
	session.Save(r, w)

	// Redirigir al usuario al formulario de inicio de sesión u a otra página
	http.Redirect(w, r, "/", http.StatusFound)
}

// Función para obtener los datos del usuario de la base de datos
func getUserData(username string) map[string]string {
	// Realizar la consulta a la base de datos para obtener los datos del usuario
	var nombre string
	var apellido string
	var email string
	err := cDB.QueryRow("SELECT nombres, apellidos, user FROM admin WHERE user = ?", username).Scan(&nombre, &apellido, &email)
	if err != nil {
		log.Println("Error al obtener los datos del usuario:", err)
		return nil
	}

	userData := map[string]string{
		"username": username,
		"nombre":   nombre,
		"apellido": apellido,
	}

	return userData
}
