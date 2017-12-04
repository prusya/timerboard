package main

import (
	"net/http"
	"html/template"
	"github.com/markbates/goth/gothic"
	"strconv"
	"encoding/json"
	"github.com/gorilla/mux"
)

var templates map[string]*template.Template

func init() {
	if templates == nil {
		templates = make(map[string]*template.Template)
	}
	templates["index"] = template.Must(
		template.New("index.html").Delims("[[", "]]").
			ParseFiles("static/index.html"))
	templates["users"] = template.Must(
		template.New("users.html").Delims("[[", "]]").
			ParseFiles("static/users.html"))
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
	err = templates["index"].ExecuteTemplate(w, "index.html", user)
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
	for key := range session.Values {
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

	var data struct {
		User  User
		Users []User
	}
	data.User = user
	data.Users = users

	err = templates["users"].ExecuteTemplate(w, "users.html", data)
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
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

func PostTimersHandler(w http.ResponseWriter, r *http.Request) {
	user, err := getUserFromSession(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !user.CanPost {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	r.ParseForm()
	var days int
	if r.Form["daysleft"][0] != "" {
		days, err = strconv.Atoi(r.Form["daysleft"][0])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		days = 0
	}
	var hours int
	if r.Form["hoursleft"][0] != "" {
		hours, err = strconv.Atoi(r.Form["hoursleft"][0])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		hours = 0
	}
	var minutes int
	if r.Form["minutesleft"][0] != "" {
		minutes, err = strconv.Atoi(r.Form["minutesleft"][0])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		minutes = 0
	}
	err = dbCreateTimer(r.Form["regioninput"][0], r.Form["systeminput"][0],
		r.Form["structureinput"][0], r.Form["rftypeinput"][0],
		r.Form["commentinput"][0], days, hours, minutes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func DeleteTimersHandler(w http.ResponseWriter, r *http.Request) {
	user, err := getUserFromSession(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !user.CanPost {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = dbDeleteTimer(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
