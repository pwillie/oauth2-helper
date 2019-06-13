package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"

	"github.com/coreos/go-oidc"
	"github.com/jpillora/opts"
	"github.com/pkg/browser"
	"golang.org/x/oauth2"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}

func main() {
	type optsConfig struct {
		Issuer       string   `opts:"help=Token issuer"`
		ClientID     string   `opts:"help=Token issuer client id"`
		ClientSecret string   `opts:"help=Token issuer client secret"`
		Scope        []string `opts:"help=Required scopes"`
	}
	cfg := optsConfig{}
	opts.Parse(&cfg)

	ctx := context.Background()
	provider, err := oidc.NewProvider(ctx, cfg.Issuer)
	if err != nil {
		log.Fatal(err)
	}
	config := oauth2.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		Endpoint:     provider.Endpoint(),
		RedirectURL:  "http://localhost:5556/auth/callback",
		Scopes:       append(cfg.Scope, oidc.ScopeOpenID),
	}

	state := randString(16)
	browser.OpenURL(config.AuthCodeURL(state))

	m := http.NewServeMux()
	s := http.Server{Addr: "localhost:5556", Handler: m}
	m.HandleFunc("/auth/callback", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("state") != state {
			http.Error(w, "state did not match", http.StatusBadRequest)
			return
		}

		oauth2Token, err := config.Exchange(ctx, r.URL.Query().Get("code"))
		if err != nil {
			http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Printf("Access Token: %s\n", oauth2Token.AccessToken)

		w.Write([]byte(`<html>`))
		w.Write([]byte("<body><h1>Success</h1><br />"))
		w.Write([]byte("<div>You may now close this window</div>"))
		w.Write([]byte("</body></head>"))

		go func() {
			if err := s.Shutdown(context.Background()); err != nil {
				log.Fatal(err)
			}
		}()
	})

	if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
