package main

//all made to work even if there is no db

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

type HTTPHandler struct {
	token string
}

func (h *HTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, h.token)
}

func main() {
	log.Println("v 1.0")
	log.Println("If you run this in production, you are opening the gates")

	weakprotection := os.Getenv("WEAKPROTECTION")
	if weakprotection == "" {
		weakprotection = "unsecure"
	}
	debug := os.Getenv("DEBUG")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.RequestURI, "favicon.ico") {
			return
		}

		if err := r.ParseForm(); err != nil {
			log.Printf("ParseForm() err: %v", err)
			fmt.Fprintf(w, "Internal error, check logs")
			return
		}

		if debug == "on" {
			log.Printf("%v", r.PostForm)
		}
		value := r.FormValue("WEAKPROTECTION")

		if value == "" || r.FormValue("WEAKPROTECTION") != weakprotection {
			log.Println("Got unvalidated request")

			if debug == "on" {
				log.Printf("%s (<<or empty) != %s", value, weakprotection)
			}

			fmt.Fprint(w, `<!doctype html>
			<html lang="en">
				<head>
				<!-- Required meta tags -->
				<meta charset="utf-8">
				<meta name="viewport" content="width=device-width, initial-scale=1">
			
				<!-- Bootstrap CSS -->
				<link href="https://cdn.jsdelivr.net/npm/bootstrap@5.1.3/dist/css/bootstrap.min.css" rel="stylesheet" integrity="sha384-1BmE4kWBq78iYhFldvKuhfTAU6auU8tT94WrHftjDbrCEXSU1oBoqyl2QvZ6jIW3" crossorigin="anonymous">
			
				<title>Get your unsecure token here</title>
				</head>
				<body onLoad="document.getElementById('token').select()">
				<div class="container">
				<div class="row">
				<div class="col">
				<h1>Don't run me in ANY Production Environment</h1>
				</div>
				</div>
				<div class="row">
				  <div class="col">
				  	<form method="post">
					  <div class="mb-3">
						  <label for="weaktoken" class="form-label">Weak Password:</label>
  						  <input type="password" class="form-control" id="WEAKPROTECTION" name="WEAKPROTECTION" placeholder="password">
					  </div>
					  <button type="submit" class="btn btn-primary">Submit</button>
					</form>
				  </div>
				</div>
				  </div>
			
				</body>
			</html>
			
			`)
		} else {
			token := "couldn't read token, make sure to inject a service account"
			tokenpath := "/var/run/secrets/kubernetes.io/serviceaccount/token"

			content, err := ioutil.ReadFile(tokenpath)
			if err == nil {
				token = string(content)
				log.Println("Serving token")
			} else {
				log.Printf("Couldn't read token on path %s", tokenpath)
				log.Println(err)
			}

			fmt.Fprintf(w, `<!doctype html>
	<html lang="en">
		<head>
		<!-- Required meta tags -->
		<meta charset="utf-8">
		<meta name="viewport" content="width=device-width, initial-scale=1">
	
		<!-- Bootstrap CSS -->
		<link href="https://cdn.jsdelivr.net/npm/bootstrap@5.1.3/dist/css/bootstrap.min.css" rel="stylesheet" integrity="sha384-1BmE4kWBq78iYhFldvKuhfTAU6auU8tT94WrHftjDbrCEXSU1oBoqyl2QvZ6jIW3" crossorigin="anonymous">
	
		<title>Get your unsecure token here</title>
		</head>
		<body onLoad="document.getElementById('token').select()">
		<div class="container">
		<div class="row">
		<div class="col">
		<h1>Don't run me in ANY Production Environment</h1>
		</div>
		</div>
		<div class="row">
		  <div class="col">
			  <div class="mb-3">
				  <label for="token" class="form-label">Get your unsecure token here</label>
				  <textarea onClick="this.select();" class="form-control" id="token" rows="10">%s</textarea>
			</div>
		  </div>
		</div>
		  </div>
	
		</body>
	</html>`, token)
		}

	})
	log.Println("Starting  on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
