package handler

import (
    "encoding/json"
    "fmt"
    "net/http"
   
    "appstore/model"
)

func uploadHandler(w http.ResponseWriter, r *http.Request) {
    // Parse from body of request to get a json object.
    fmt.Println("Received one upload request")
    decoder := json.NewDecoder(r.Body)
    var app model.App
    if err := decoder.Decode(&app); err != nil {
        panic(err)
    }

    fmt.Fprintf(w, "Upload request received: %s\n", app.Description)
}
