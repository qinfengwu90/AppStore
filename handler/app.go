package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
    jwt "github.com/form3tech-oss/jwt-go"
	"appstore/model"
	"appstore/service"
    "github.com/gorilla/mux"  
	"github.com/pborman/uuid"
)

func uploadHandler(w http.ResponseWriter, r *http.Request) {
    // Parse from body of request to get a json object.
    fmt.Println("Received one upload request")
    
    user := r.Context().Value("user")
    claims := user.(*jwt.Token).Claims
    username := claims.(jwt.MapClaims)["username"]

    app := model.App{
        Id: uuid.New(),
        User: username.(string),
        Title: r.FormValue("title"),
        Description: r.FormValue("description"),
    }

    price, err := strconv.ParseFloat(r.FormValue("price"), 64)
    fmt.Printf("%v,%T", price, price)
    if err != nil {
        fmt.Println(err)
    }
    app.Price = int(price * 100.0)

    file, _, err := r.FormFile("media_file")
    if err != nil {
        http.Error(w, "Media file is unavailble", http.StatusBadRequest)
        fmt.Printf("Media file is not available %v\n", err)
        return
    }

    err = service.SaveApp(&app, file)
    if err != nil {
        http.Error(w, "Failed to save app to database", http.StatusBadRequest)
        fmt.Printf("Failed to save app to database %v\n", err)
        return
    }
    fmt.Println("App is saved successfully.")
    fmt.Fprintf(w, "App is saved successfully: %s\n", app.Description)
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Println("Received one search request")
    w.Header().Set("Content-Type", "aplication/json")
    title := r.URL.Query().Get("title")
    description := r.URL.Query().Get("description")

    var apps []model.App
    var err error
    apps, err = service.SearchApps(title, description)
    if err != nil {
        http.Error(w, "Failed to read apps from the backend", http.StatusInternalServerError)
        return
    }

    js, err := json.Marshal(apps)
    if err != nil {
        http.Error(w, "Failed to parse apps into JSON format", http.StatusInternalServerError)
        return
    }
    w.Write(js)
}

func checkoutHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Println("Received one checkout request")
    w.Header().Set("Content-Type", "text/plain")
    appID := r.FormValue("appID")
    s, err := service.CheckoutApp(r.Header.Get("Origin"), appID)
    if err != nil {
        fmt.Println("Checkout failed")
        w.Write([]byte(err.Error()))
        return
    }
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(s.URL))
    fmt.Println("Checkout process started!")
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Println("Received one request for delete")

    user := r.Context().Value("user")
    claims := user.(*jwt.Token).Claims
    username := claims.(jwt.MapClaims)["username"].(string)
    id := mux.Vars(r)["id"]

    if err := service.DeleteApp(id, username); err != nil {
        http.Error(w, "Failed to delete app from backend", http.StatusInternalServerError)
        fmt.Printf("Failed to delete app from backend %v\n", err)
        return
    }
    fmt.Println("App is deleted successfully")
}