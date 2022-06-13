package main

import (
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

func main() {

	// load variables in the file to environment
	err := godotenv.Load("local.env")
	if err != nil {
		log.Fatalf("Error: %s", err)
	}

	// get the initial variables from environment
	accessToken := os.Getenv("ACCESS_TOKEN")
	authToken := os.Getenv("AUTH_TOKEN")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		_, err := w.Write([]byte("Service is running."))
		if err != nil {
			log.Fatalf("Error: %s", err)
			return
		}
	})

	http.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {

		switch r.Method {

		// send current access token back
		case "GET":
			w.WriteHeader(200)
			_, err := w.Write([]byte(accessToken))
			if err != nil {
				log.Fatalf("Error: %s", err)
				return
			}

		// set new access token
		case "POST":
			err := r.ParseForm()
			if err != nil {
				log.Fatalf("Error: %s", err)
				return
			}
			token := r.FormValue("token")
			auth := r.FormValue("auth")

			// verify request data (client need to send auth which matches the AUTH_TOKEN in environment)
			if token == "" || auth != authToken {
				w.WriteHeader(200)
				_, err = w.Write([]byte("Invalid data."))
				if err != nil {
					log.Fatalf("Error: %s", err)
					return
				}
			} else {
				accessToken = token
				w.WriteHeader(200)
				_, err = w.Write([]byte(accessToken))
				if err != nil {
					log.Fatalf("Error: %s", err)
					return
				}
			}
		}
	})

	log.Fatal(http.ListenAndServe(":8080", nil))

}
