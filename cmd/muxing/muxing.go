package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"strconv"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

/**
Please note Start functions is a placeholder for you to start your own solution.
Feel free to drop gorilla.mux if you want and use any other solution available.

main function reads host/port from env just for an example, flavor it following your taste
*/

// Start /** Starts the web server listener on given host and port.
func Start(host string, port int) {
	mux := mux.NewRouter()
	//mux := http.NewServeMux()

	reqCtx := &reqContext{verbose: true}

	mux.Handle("/bad", handlers.LoggingHandler(os.Stdout,
		Handler(reqCtx, "GET", http.HandlerFunc(getBad)))).Methods("GET")

	mux.Handle("/name/{param}", handlers.LoggingHandler(os.Stdout,
		Handler(reqCtx, "GET", http.HandlerFunc(getName))))

	mux.Handle("/data", handlers.LoggingHandler(os.Stdout,
		Handler(reqCtx, "POST", http.HandlerFunc(getData)))).Methods("POST")

	mux.Handle("/headers", handlers.LoggingHandler(os.Stdout,
		Handler(reqCtx, "POST", http.HandlerFunc(getHeaders))))

	mux.Handle("/get-echo", handlers.LoggingHandler(os.Stdout,
		Handler(reqCtx, "GET", http.HandlerFunc(getEcho))))

	log.Println(fmt.Printf("Starting API server on %s:%d\n", host, port))
	if err := http.ListenAndServe(fmt.Sprintf("%s:%d", host, port), mux); err != nil {
		log.Fatal(err)
	}
}

//main /** starts program, gets HOST:PORT param and calls Start func.
func main() {
	host := os.Getenv("HOST")
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		port = 8081
	}
	Start(host, port)
}

// getEcho Retrieves the echo response
func getEcho(w http.ResponseWriter, r *http.Request) {
	log.Println("Echoing back request made to " + r.URL.Path + " to client (" + r.RemoteAddr + ")")

	w.Header().Set("Access-Control-Allow-Origin", "*")

	// allow pre-flight headers
	w.Header().Set("Access-Control-Allow-Headers", "Content-Range, Content-Disposition, Content-Type, ETag")

	err := r.Write(w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func getBad(w http.ResponseWriter, r *http.Request) {
	httpError(w, http.StatusInternalServerError) // 500
}

func getName(w http.ResponseWriter, r *http.Request) {
	//urlTokens := strings.Split(r.RequestURI, "/")
	//param := fmt.Sprintf("Hello, %s!", urlTokens[len(urlTokens)-1])
	vars := mux.Vars(r)
	param := fmt.Sprintf("Hello, %s!", vars["param"])
	_, err := w.Write([]byte(param))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func getData(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	param := fmt.Sprintf("I got message:\n%s", string(body))
	_, err = w.Write([]byte(param))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func getHeaders(w http.ResponseWriter, r *http.Request) {
	a := r.Header["a"]
	b := r.Header["b"]
	ai, err := strconv.Atoi(a[0])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	bi, err := strconv.Atoi(b[0])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("a+b", strconv.Itoa(ai+bi))
}

func Handler(reqCtx *reqContext, method string, h http.Handler) http.Handler {
	return handler{reqCtx, method, h}
}

type reqContext struct {
	verbose bool
}

type handler struct {
	reqCtx  *reqContext
	method  string
	handler http.Handler
}

func (vh handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if vh.reqCtx.verbose {
		reqDump, err := httputil.DumpRequest(r, true)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Printf("\n%s", reqDump)
	}
	if r.Method != vh.method {
		httpError(w, http.StatusMethodNotAllowed) // 405
		return
	}
	vh.handler.ServeHTTP(w, r)
}

func httpError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}
