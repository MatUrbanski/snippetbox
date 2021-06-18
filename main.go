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
  // Check if the current request URL path exactly matches "/". If it doesn't, use
  // the http.NotFound() function to send a 404 response to the client.
  // Importantly, we then return from the handler. If we don't return, the handler
  // would keep execurint and also write the "Hello from Snippetbox".
  if r.URL.Path != "/" {
    http.NotFound(w, r)
    return
  }

  w.Write([]byte("Hello from SnippetBox"))
}

// Add a showSnippet handler function.
func showSnippet(w http.ResponseWriter, r *http.Request) {
  w.Write([]byte("Display a specific snippet..."))
}

// Add a createSnippet handler function.
func createSnippet(w http.ResponseWriter, r *http.Request) {
  // Use r.Method to check whether the request is using POST or not. Note that
  // http.MethodPost is a constant equal to the string "POST".
  if r.Method != http.MethodPost {
    // Use the Header().Set() method to add an 'Allow: POST' header to the
    // response header map. The first paramter is the header name, and
    // the second parameter is the header value.
    w.Header().Set("Allow", http.MethodPost)

    // Use the http.Error() function to send a 405 status code and "Method not
    // Allowed" string as the response body.
    http.Error(w, "Method Not Allowed", 405)
    return
  }

  w.Write([]byte("Create a specific snippet..."))
}

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
