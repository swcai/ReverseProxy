package main 

import (
    "log"
    "flag"
    "fmt"
    "time"
    "net/http"
    "math/rand"
)

const APP_VERSION = "0.1"

// The flag package provides a default help printer via -h switch
var versionFlag *bool = flag.Bool("v", false, "Print the version number.")
var port *int = flag.Int("p", 8081, "The port for web server.")

func randomHandler(w http.ResponseWriter, r *http.Request) {
    output := fmt.Sprintf("%v", rand.Int())
    w.Write([]byte(output))
    log.Printf("done %v %v", *port, output)
}

func main() {
    flag.Parse() // Scan the arguments list 

    if *versionFlag {
        fmt.Println("Version:", APP_VERSION)
    }
    
    rand.Seed(time.Now().UnixNano())
    http.HandleFunc("/random", randomHandler)
    addr := fmt.Sprintf("localhost:%d", *port)
    http.ListenAndServe(addr, nil)
}

