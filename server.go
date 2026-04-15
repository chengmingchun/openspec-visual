package main

import (
	"encoding/json"
	"net/http"
)

func StartLocalServer(svc *OpenSpecService) {
	mux := http.NewServeMux()
	
	mux.HandleFunc("/api/config", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == "OPTIONS" {
			w.WriteHeader(200)
			return
		}

		if r.Method == "GET" {
			cfg := svc.GetConfig()
			json.NewEncoder(w).Encode(cfg)
			return
		}
		if r.Method == "POST" {
			var p struct{ APIKey, BaseURL, Model string }
			if err := json.NewDecoder(r.Body).Decode(&p); err == nil {
				svc.SaveConfig(p.APIKey, p.BaseURL, p.Model)
				w.WriteHeader(200)
			}
		}
	})

	mux.HandleFunc("/api/generate", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == "OPTIONS" { w.WriteHeader(200); return }

		var p struct{ FeatureName, Content string }
		json.NewDecoder(r.Body).Decode(&p)
		err := svc.GenerateOpenSpecStructure(p.FeatureName, p.Content)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		w.WriteHeader(200)
	})

	mux.HandleFunc("/api/list", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		node, _ := svc.ListOpenSpecFiles()
		json.NewEncoder(w).Encode(node)
	})

	mux.HandleFunc("/api/read", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		path := r.URL.Query().Get("path")
		data, err := svc.ReadFileContent(path)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		w.Write([]byte(data))
	})

	mux.HandleFunc("/api/prompt", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == "OPTIONS" { w.WriteHeader(200); return }

		var p struct{ Prompt, System string }
		json.NewDecoder(r.Body).Decode(&p)
		res, err := svc.RunPrompt(p.Prompt, p.System)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"result": res})
	})

	go http.ListenAndServe("127.0.0.1:38192", mux)
}
