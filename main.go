package main

import (
	"errors"
	"log"
	"net/http"
	"os"
	"strings"

	"git.xeserv.us/ponychat/shoutpage/atheme"
	"github.com/codegangsta/negroni"
	"github.com/yosssi/ace"
)

func init() {
	if os.Getenv("ATHEME_URL") == "" {
		panic("Need ATHEME_URL")
	}

	if os.Getenv("BNC_URL") == "" {
		panic("Need BNC_URL")
	}

	if os.Getenv("OUTPUT_PATH") == "" {
		panic("Need OUTPUT_PATH")
	}
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		doTemplate("index", rw, r, nil)
	})

	mux.HandleFunc("/signup", func(rw http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			handleError(rw, r, errors.New("wrong method"))
			return
		}

		err := r.ParseForm()
		if err != nil {
			handleError(rw, r, err)
			return
		}

		username := r.Form.Get("username")
		password := r.Form.Get("password")

		a, _ := atheme.NewAtheme(os.Getenv("ATHEME_URL"))
		err = a.Login(username, password)
		if err != nil {
			handleError(rw, r, err)
			return
		}

		info, err := a.NickServ.OwnInfo()
		if err != nil {
			handleError(rw, r, err)
			return
		}

		groupsRaw, ok := info["groups"]
		if !ok {
			doTemplate("notinprogram", rw, r, username)
			return
		}

		groups := strings.Split(groupsRaw, ", ")

		for _, group := range groups {
			if group == "!bncusers" {
				goto ok
			}
		}

		doTemplate("notinprogram", rw, r, username)
		return

	ok:
		MakeUser(username, password)

		doTemplate("success", rw, r, struct {
			Username string
			URL      string
		}{
			Username: username,
			URL:      os.Getenv("BNC_URL"),
		})
	})

	n := negroni.Classic()

	n.UseHandler(mux)
	n.Run(":3000")
}

func handleError(rw http.ResponseWriter, r *http.Request, err error) {
	rw.WriteHeader(500)

	data := struct {
		Path   string
		Reason string
	}{
		Path:   r.URL.String(),
		Reason: err.Error(),
	}

	log.Printf("%#v", err)

	tpl, err := ace.Load("layout", "error", nil)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tpl.Execute(rw, data); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}

func doTemplate(name string, rw http.ResponseWriter, r *http.Request, data interface{}) {
	tpl, err := ace.Load("layout", name, nil)
	if err != nil {
		handleError(rw, r, err)
		return
	}

	if err := tpl.Execute(rw, data); err != nil {
		handleError(rw, r, err)
		return
	}
}
