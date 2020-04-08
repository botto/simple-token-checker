package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/caarlos0/env"
	"github.com/gobuffalo/packr/v2"
	"github.com/hashicorp/vault/api"
)

var vault *api.Client
var indexTpl string
var errorTpl string

func main() {
	var httpClient *http.Client
	var err error
	cfg := config{}
	if err = env.Parse(&cfg); err != nil {
		log.Fatalf("Failed to parse config %s", err)
	}
	handleSignals()
	if len(cfg.CAFilePath) > 0 {
		httpClient, err = ownCAHttpClient(cfg.CAFilePath)
		if err != nil {
			log.Fatalf("Could not create custom http client: %s", err)
		}
	}

	vault, err = api.NewClient(&api.Config{
		Address:    cfg.VaultAddr,
		HttpClient: httpClient, // If nil, it uses a defult vault api
	})
	if err != nil {
		log.Fatalf("Problem with vault client init: %s", err)
	}
	vault.SetToken(cfg.VaultToken)
	box := packr.New("box", "./assets")
	indexTpl, err = box.FindString("index.html")
	if err != nil {
		log.Fatal("Could not parse index.html")
	}
	errorTpl, err = box.FindString("error.html")
	if err != nil {
		log.Fatal("Could not parse error.html")
	}
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(box)))
	http.HandleFunc("/", index)
	listenDetails := fmt.Sprintf("%s:%d", cfg.ListenAddr, cfg.ListenPort)
	err = http.ListenAndServe(listenDetails, nil)
	if err != nil {
		log.Fatalf("Problem starting server: %s", err)
	}
}

// If a user is redirected to the login service.
func index(w http.ResponseWriter, r *http.Request) {
	var tokenString string

	tokenCookie, _ := r.Cookie("X-Vault-Token")
	if tokenCookie == nil || len(tokenCookie.Value) <= 0 {
		tokenHeader, ok := r.Header["X-Vault-Token"]
		if ok && len(tokenHeader[0]) > 0 {
			tokenString = tokenHeader[0]
		}
	} else {
		tokenString = tokenCookie.Value
	}

	if len(tokenString) > 0 {
		_, err := vault.Auth().Token().Lookup(tokenString)
		if err == nil {
			t, _ := template.New("index").Parse(indexTpl)
			t.Execute(w, struct {
				Token string
			}{
				Token: tokenString,
			})
			return
		}
		t, _ := template.New("error").Parse(errorTpl)
		t.Execute(w, struct {
			Token string
		}{
			Token: tokenString,
		})
		return
	}
	t, _ := template.New("error").Parse(errorTpl)
	t.Execute(w, struct{ Token string }{})
}
