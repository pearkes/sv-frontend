package main

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/pearkes/Dropbox-Go/dropbox"
	"github.com/pearkes/sv-frontend/stats"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"
)

var templates = template.Must(template.ParseFiles("files/index.html", "files/error.html", "files/cancel.html", "files/setup.html", "files/help.html"))

type Page struct {
	Title string
	Uid   string
}

func userPageHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s %s", r.Method, r.URL, r.RemoteAddr)
	u, err := db.FindByName(r.Host)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	conn := red.Redis.Get()
	defer conn.Close()

	keyName := fmt.Sprintf("page:%v:%s", u.Id, u.DropboxUid)
	page, err := redis.String(conn.Do("GET", keyName))
	if err != nil {
		if err.Error() == "redigo: nil returned" {
			setupPage(w, r)
			return
		}
		log.Printf("Error retreiving page: %s", err.Error())
		errorPage(w, r, err)
		return
	}
	w.Write([]byte(page))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s %s", r.Method, r.URL, r.RemoteAddr)
	p := &Page{Title: "homepage!"}
	templates.ExecuteTemplate(w, "index.html", p)
}

func authCallbackHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s %s", r.Method, r.URL, r.RemoteAddr)

	code := r.FormValue("code")

	s := dropbox.Session{
		AppKey:      DROPBOX_KEY,
		AppSecret:   DROPBOX_SECRET,
		AccessType:  "app_folder",
		RedirectUri: DROPBOX_CALLBACK,
	}

	token, uid, err := s.ObtainToken(code)

	if err != nil {
		log.Printf("Error obtaining token: %s", err.Error())
		errorPage(w, r, err)
		return
	}

	if r.FormValue("error") == "access_denied" {
		templates.ExecuteTemplate(w, "cancel.html", "")
		metrics.Event(stats.CANCEL_DBX_AUTH)
		return
	}

	u, err := db.FindByUid(uid)

	if u.DropboxUid == uid {
		err = db.UpdateToken(u, token)
	} else {
		err = db.NewUser(token, uid)
	}

	if err != nil {
		log.Printf("Error storing user token: %s", err.Error())
		errorPage(w, r, err)
		return
	}

	metrics.Event(stats.ACCOUNT_CREATE)

	// Create the default dropbox
	go fillDropbox(token)

	// redirect them to their new page
	url := fmt.Sprintf("http://%s.%s", uid, host)

	if u.DropboxUid == uid {
		url = fmt.Sprintf("http://%s", u.Name)
	}

	http.Redirect(w, r, url, 302)
}

func authInitHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s %s", r.Method, r.URL, r.RemoteAddr)
	http.Redirect(w, r, newDropboxUrl(), 302)
}

func errorPage(w http.ResponseWriter, r *http.Request, err error) {
	metrics.Event(stats.RENDER_ERROR)
	code := strconv.FormatInt(time.Now().UTC().Unix(), 10)
	message := fmt.Sprintf("%s", code)
	log.Printf("Logging user facing error: %s, %s", message, err.Error())
	templates.ExecuteTemplate(w, "error.html", message)
}

func setupPage(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "setup.html", r.Host)
}

func helpHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "help.html", r.Host)
}
