package main

import (
  "log"
  "net/http"
)

// Define a home handler function which writes a byte slice containing
// "Hello from Snippetbox" as the reponse body.
// The http.ResponseWriter parameter provides methods for assembling a HTTP
// response and sending it to the user, and the *http.Request parameter is a
// pointer to a struct which holds information about the current request
// (like the HTTP method and the URL being requested).
func home(w http.ResponseWriter, r *http.Request) {
  w.Write([]byte("Hello from Snippetbox"))
}

func main() {
  // Use the http.newServeMux() function to initialize a new servemux, then
  // register the home function as the handler for the "/" URL pattern.
  mux := http.NewServeMux()
  mux.HandleFunc("/", home)

  // Use the http.ListenAndServe() function to start a new web server. We pass in
  // two parameters: the TCP network addres to listen on (in this case ":4000")
  // and the servemux we just created. If http.ListenAndServe() returns an error
  // we use log.Fatal() function to log the error message and exit. Note
  // that any error returned by http.ListenAndServe() is always non-nil.
  log.Println("Starting server on :4000")
  err := http.ListenAndServe(":4000", mux)
  log.Fatal(err)
}
