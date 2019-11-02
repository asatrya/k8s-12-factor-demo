package main

import (
    "fmt"
    "net/http"
    "os"
    "time"

    "github.com/gorilla/handlers"
    "github.com/gorilla/mux"
)

var started = time.Now()

// used to dump headers for debugging
func indexHandler(w http.ResponseWriter, r *http.Request) {

    startTime := time.Now()

    // disable cache
    w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
    w.Header().Set("Pragma", "no-cache")
    w.Header().Set("Expires", "0")

    // set hostname (used for demo)
    hostname, err := os.Hostname()
    if err != nil {
        fmt.Fprint(w, "Error:", err)
    }

    fmt.Fprintf(w, "Service/Pod Info: %s\n\n", os.Getenv("ENV"))
    fmt.Fprintf(w, "Served-By: %v\n", hostname)
    fmt.Fprintf(w, "Serving-Time: %s\n", time.Now().Sub(startTime))
    
    fmt.Fprintf(w, "\nEnvironment variables: %s\n\n", os.Getenv("ENV"))
    fmt.Fprintf(w, "ENV: %s\n", os.Getenv("ENV"))
    fmt.Fprintf(w, "DB_HOST: %s\n", os.Getenv("DB_HOST"))
    fmt.Fprintf(w, "DB_PORT: %s\n", os.Getenv("DB_PORT"))
    fmt.Fprintf(w, "DB_USER: %s\n", os.Getenv("DB_USER"))
    fmt.Fprintf(w, "DB_PASSWORD: %s\n", os.Getenv("DB_PASSWORD"))
    return

}

func healthzHandler(w http.ResponseWriter, r *http.Request) {
    duration := time.Now().Sub(started)
    if duration.Seconds() > 10 {
        http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
    } else {
        fmt.Fprintf(w, "ok: %v", duration.Seconds())
    }
    return
}

// mux
var router = mux.NewRouter()

func main() {

    router.HandleFunc("/", indexHandler)
    router.HandleFunc("/healthz", healthzHandler)
    http.Handle("/", router)

    fmt.Println("Listening on port 5005...")
    http.ListenAndServe(":5005", handlers.CompressHandler(router))

}
