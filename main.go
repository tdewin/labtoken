package main

//all made to work even if there is no db

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
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

	token := "couldn't read token, make sure to inject a service account"
	tokenpath := "/var/run/secrets/kubernetes.io/serviceaccount"

	content, err := ioutil.ReadFile(tokenpath)
	if err == nil {
		token = string(content)
	} else {
		log.Printf("Couldnt read token on path %s", tokenpath)
		log.Println(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
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
	})
	log.Fatal(http.ListenAndServe(":8080", nil))

}
