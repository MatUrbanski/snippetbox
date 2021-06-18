package main

import (
  "log"
  "net/http"
)

func main() {
  // Use the http.newServeMux() function to initialize a new servemux.
  mux := http.NewServeMux()

  // Register the home function as the handler for the "/" URL pattern.
  mux.HandleFunc("/", home)

  // Register the showSnippet function as the handler for the "/snippet" URL pattern.
  mux.HandleFunc("/snippet", showSnippet)

  // Register the createSnippet function as the handler for the "/snippet/create" URL pattern.
  mux.HandleFunc("/snippet/create", createSnippet)

  // Use the http.ListenAndServe() function to start a new web server. We pass in
  // two parameters: the TCP network addres to listen on (in this case ":4000")
  // and the servemux we just created. If http.ListenAndServe() returns an error
  // we use log.Fatal() function to log the error message and exit. Note
  // that any error returned by http.ListenAndServe() is always non-nil.
  log.Println("Starting server on :4000")
  err := http.ListenAndServe(":4000", mux)
  log.Fatal(err)
}
