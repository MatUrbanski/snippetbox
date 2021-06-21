package main

import(
  "net/http"
  "github.com/justinas/alice"
)

func(app *application) routes() http.Handler {
  // Create a middleware chain containing our 'standard' middleware
  // which will be used for every request our application receives.
  standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

  // Use the http.newServeMux() function to initialize a new servemux.
  mux := http.NewServeMux()

  // Register the home function as the handler for the "/" URL pattern.
  mux.HandleFunc("/", app.home)

  // Register the showSnippet function as the handler for the "/snippet" URL pattern.
  mux.HandleFunc("/snippet", app.showSnippet)

  // Register the createSnippet function as the handler for the "/snippet/create" URL pattern.
  mux.HandleFunc("/snippet/create", app.createSnippet)

  // Create a file server which serves files out of the "./ui/static" directory.
  // Note that the path given to the http.Dir function is relative to the project
  // directory root.
  fileServer := http.FileServer(http.Dir("./ui/static"))

  // Use the mux.Handle() function to register the file server as the handler for
  // all URL paths that start with "/static/". For matching paths, we strip the
  // "/static" prefix before the request reaches the file server.
  mux.Handle("/static/", http.StripPrefix("/static", fileServer))

  // Return the 'standard' middleware chain followed by the servemux.
  return standardMiddleware.Then(mux)
}
