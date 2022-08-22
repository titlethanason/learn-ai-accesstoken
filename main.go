package main

import (
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"strconv"
)

type Resp struct {
	Error   int    `json:"error"`
	Message string `json:"message"`
}

func sendEmail(accessToken string) {
	from := "thanason.eiam@mail.kmutt.ac.th"
	password := os.Getenv("EMAIL_PASSWORD")

	// Receiver email address.
	to := []string{
		"titlethanason@gmail.com",
		"thanason.e@mail.kmutt.ac.th",
	}

	// smtp server configuration.
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// Message.
	tmpMessage := "From: " + from + "\r\n" + "To: "
	for _, s := range to {
		tmpMessage = tmpMessage + "," + s
	}
	tmpMessage = tmpMessage + "\r\nSubject: Facebook conversion API error\r\n\r\n" + "Error sending request to Facebook conversion API. This might be a problem about invalid token: " + accessToken + "\r\n"
	message := []byte(tmpMessage)

	// Authentication.
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Sending email.
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if err != nil {
		fmt.Println(err)
		return
	}
	log.Println("Email sent successfully.")
}

func main() {

	// load variables in the file to environment
	err := godotenv.Load("local.env")
	if err != nil {
		log.Fatalf("Error: %s", err)
	}

	// get the initial variables from environment
	accessToken := os.Getenv("ACCESS_TOKEN")
	authToken := os.Getenv("AUTH_TOKEN")
	sentEmail := false

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
			_, err := w.Write([]byte(".... " + accessToken[len(accessToken)-10:]))
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
				resp := &Resp{Error: 1, Message: "Invalid data"}
				jsonResp, err := json.Marshal(resp)
				if err != nil {
					log.Fatalf("Error at json marshal: %s", err)
				}
				w.WriteHeader(200)
				_, err = w.Write(jsonResp)
				if err != nil {
					log.Fatalf("Error: %s", err)
					return
				}
			} else {
				accessToken = token
				sentEmail = false
				resp := &Resp{Error: 0, Message: accessToken}
				jsonResp, err := json.Marshal(resp)
				if err != nil {
					log.Fatalf("Error at json marshal: %s", err)
				}
				w.WriteHeader(200)
				_, err = w.Write(jsonResp)
				if err != nil {
					log.Fatalf("Error: %s", err)
					return
				}
			}
		}
	})

	http.HandleFunc("/new", func(w http.ResponseWriter, r *http.Request) {

		switch r.Method {

		// return email sending status
		case "GET":
			w.WriteHeader(200)
			_, err := w.Write([]byte("Email has sent: " + strconv.FormatBool(sentEmail)))
			if err != nil {
				log.Fatalf("Error: %s", err)
				return
			}

		// send email and wait for new token
		case "POST":
			err := r.ParseForm()
			if err != nil {
				log.Fatalf("Error: %s", err)
				return
			}
			token := r.FormValue("token")

			if token == accessToken {

				// send email
				if !sentEmail {
					sendEmail(accessToken)
					sentEmail = true
				}

				// wait for new token
				resp := &Resp{Error: 1, Message: "Waiting for new token"}
				jsonResp, err := json.Marshal(resp)
				if err != nil {
					log.Fatalf("Error at json marshal: %s", err)
				}
				w.WriteHeader(200)
				_, err = w.Write(jsonResp)
				if err != nil {
					log.Fatalf("Error: %s", err)
					return
				}
			} else {
				// send new token back
				resp := &Resp{Error: 0, Message: accessToken}
				jsonResp, err := json.Marshal(resp)
				if err != nil {
					log.Fatalf("Error at json marshal: %s", err)
				}
				w.WriteHeader(200)
				_, err = w.Write(jsonResp)
				if err != nil {
					log.Fatalf("Error: %s", err)
					return
				}
			}
		}
	})

	port := os.Getenv("HOST_PORT")
	if port == "" {
		port = "80"
	}
	log.Fatal(http.ListenAndServe(":"+port, nil))

}
