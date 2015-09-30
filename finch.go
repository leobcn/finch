package main // import "github.com/thraxil/finch"

import (
	_ "expvar"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
)

var templateDir = "templates"

func makeHandler(fn func(http.ResponseWriter, *http.Request, *site), s *site) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fn(w, r, s)
	}
}

func main() {
	p := newPersistence(os.Getenv("FINCH_DB_FILE"))
	defer p.Close()
	templateDir = os.Getenv("FINCH_TEMPLATE_DIR")
	s := newSite(
		p,
		os.Getenv("FINCH_BASE_URL"),
		sessions.NewCookieStore([]byte(os.Getenv("FINCH_SECRET"))),
		os.Getenv("FINCH_ITEMS_PER_PAGE"))

	http.HandleFunc("/", makeHandler(indexHandler, s))
	http.HandleFunc("/post/", makeHandler(postHandler, s))
	http.HandleFunc("/search/", makeHandler(searchHandler, s))

	http.HandleFunc("/u/", makeHandler(userDispatch, s))

	// authy stuff
	http.HandleFunc("/register/", makeHandler(registerHandler, s))
	http.HandleFunc("/login/", makeHandler(loginHandler, s))
	http.HandleFunc("/logout/", makeHandler(logoutHandler, s))

	// static misc.
	http.HandleFunc("/favicon.ico", faviconHandler)
	http.Handle("/media/", http.StripPrefix("/media/",
		http.FileServer(http.Dir(os.Getenv("FINCH_MEDIA_DIR")))))
	log.Println("running on " + os.Getenv("FINCH_PORT"))
	log.Fatal(http.ListenAndServe(":"+os.Getenv("FINCH_PORT"), nil))
}
