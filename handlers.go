package main

import (
	"net/http"
	"html/template"
	"github.com/markbates/goth/gothic"
	"os/user"
)

var templates map[string]*template.Template

func init() {
	if templates == nil {
		templates = make(map[string]*template.Template)
	}
	templates["index"] = template.Must(template.ParseFiles("static/index.html"))
	templates["users"] = template.Must(template.ParseFiles("static/users.html"))
}

func getUserFromSession(r *http.Request) (User, error) {
	session, _ := cookieStore.Get(r, "auth")
	var user User
	var err error
	if name, ok := session.Values["name"]; ok {
		err = db.One("Name", name, &user)
	} else {
		user = User{}
	}
	return user, err
}

func GetIndexHandler(w http.ResponseWriter, r *http.Request) {
	user, err := getUserFromSession(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = templates["index"].ExecuteTemplate(w, "index", user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func GetEveCallbackHandler(w http.ResponseWriter, r *http.Request) {
	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	session, _ := cookieStore.Get(r, "auth")
	session.Values["name"] = user.NickName
	session.Save(r, w)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func GetLogoutHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := cookieStore.Get(r, "auth")
	for key, _ := range session.Values {
		delete(session.Values, key)
	}
	session.Save(r, w)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	user, err := getUserFromSession(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !user.IsAdmin {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	users, err := dbGetUsers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = templates["users"].ExecuteTemplate(w, "users", users)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func PostUsersHandler(w http.ResponseWriter, r *http.Request) {
	user, err := getUserFromSession(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !user.IsAdmin {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	r.ParseForm()
}

func GetTimersHandler(w http.ResponseWriter, r *http.Request) {
	user, err := getUserFromSession(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !user.CanRead {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	timers, err := dbGetTimers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	j, err := json.Marshal(timers)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func PostTimersHandler(w http.ResponseWriter, r *http.Request) {}

func DeleteTimersHandler(w http.ResponseWriter, r *http.Request) {}
