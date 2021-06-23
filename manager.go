package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
)

// Manager holds a list of the crypto and tokens we are watching
type Manager struct {
	Watching map[string]*Token
	Cache    *redis.Client
	Context  context.Context
	sync.RWMutex
}

// NewManager stores all the information about the current tokens being watched and
// listens for api requests on 8080
func NewManager(address string) *Manager {
	m := &Manager{
		Watching: make(map[string]*Token),
	}

	r := mux.NewRouter()
	r.HandleFunc("/ticker", m.AddToken).Methods("POST")
	r.HandleFunc("/ticker/{id}", m.DeleteToken).Methods("DELETE")
	r.HandleFunc("/ticker", m.GetTokens).Methods("GET")

	srv := &http.Server{
		Addr:         address,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}

	logger.Infof("Starting api server on %s...", address)

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	return m
}

// TokenRequest represents the json coming in from the request
type TokenRequest struct {
	Contract  string `json:"contract"`
	Token     string `json:"discord_bot_token"`
	Name      string `json:"name"`
	Nickname  bool   `json:"set_nickname"`
	Frequency int    `json:"frequency" default:"60"`
	Currency  string `json:"currency" default:"0x2791Bca1f2de4661ED88A30C99A7a9449Aa84174"`
}

// AddToken adds a new Token or crypto to the list of what to watch
func (m *Manager) AddToken(w http.ResponseWriter, r *http.Request) {
	m.Lock()
	defer m.Unlock()

	logger.Debugf("Got an API request to add a ticker")

	// read body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.Errorf("Error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error reading body: %v", err)
		return
	}

	// unmarshal into struct
	var tokenReq TokenRequest
	if err := json.Unmarshal(body, &tokenReq); err != nil {
		logger.Errorf("Error unmarshalling: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error unmarshalling: %v", err)
		return
	}

	// ensure token is set
	if tokenReq.Token == "" {
		logger.Error("Discord token required")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Error: token required")
		return
	}

	// ensure currency is set
	if tokenReq.Currency == "" {
		tokenReq.Currency = "0x2791Bca1f2de4661ED88A30C99A7a9449Aa84174"
	}

	// ensure freq is set
	if tokenReq.Frequency == 0 {
		tokenReq.Frequency = 60
	}

	// ensure name is set
	if tokenReq.Name == "" {
		logger.Error("Name required for token")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Error: Name required")
		return
	}

	// check if already existing
	if _, ok := m.Watching[strings.ToUpper(tokenReq.Contract)]; ok {
		logger.Error("Error: ticker already exists")
		w.WriteHeader(http.StatusConflict)
		return
	}

	token := NewToken(tokenReq.Contract, tokenReq.Token, tokenReq.Name, tokenReq.Nickname, tokenReq.Frequency, tokenReq.Currency)
	m.addToken(token)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(token)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	return
}

func (m *Manager) addToken(token *Token) {
	m.Watching[token.Contract] = token
}

// DeleteToken addds a new token or crypto to the list of what to watch
func (m *Manager) DeleteToken(w http.ResponseWriter, r *http.Request) {
	m.Lock()
	defer m.Unlock()

	logger.Debugf("Got an API request to delete a ticker")

	vars := mux.Vars(r)
	id := strings.ToUpper(vars["id"])

	if _, ok := m.Watching[id]; !ok {
		logger.Error("Error: no ticker found")
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "Error: ticker not found")
		return
	}
	// send shutdown sign
	m.Watching[id].Shutdown()

	// remove from cache
	delete(m.Watching, id)

	logger.Infof("Deleted ticker %s", id)
	w.WriteHeader(http.StatusNoContent)
}

// GetTokens returns a list of what the manager is watching
func (m *Manager) GetTokens(w http.ResponseWriter, r *http.Request) {
	m.RLock()
	defer m.RUnlock()
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(m.Watching); err != nil {
		logger.Errorf("Error serving request: %v", err)
		fmt.Fprintf(w, "Error: %v", err)
	}
}
