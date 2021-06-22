package main

import(
  "net/http"
  "github.com/bmizerany/pat"
  "github.com/justinas/alice"
)

func(app *application) routes() http.Handler {
  // Create a middleware chain containing our 'standard' middleware
  // which will be used for every request our application receives.
  standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

  // Use the pat.New() function to initialize a new Pat router.
  mux := pat.New()

  // Register the home function as the handler for the "/" URL pattern.
  mux.Get("/", http.HandlerFunc(app.home))

  // Register the createSnippetForm function as the handler for the GET "/snippet/create" URL pattern.
  mux.Get("/snippet/create", http.HandlerFunc(app.createSnippetForm))

  // Register the createSnippet function as the handler for the POST "/snippet/create" URL pattern.
  mux.Post("/snippet/create", http.HandlerFunc(app.createSnippet))

  // Register the showSnippet function as the handler for the "/snippet/:id" URL pattern.
  mux.Get("/snippet/:id", http.HandlerFunc(app.showSnippet))

  // Create a file server which serves files out of the "./ui/static" directory.
  // Note that the path given to the http.Dir function is relative to the project
  // directory root.
  fileServer := http.FileServer(http.Dir("./ui/static"))

  // Use the mux.Handle() function to register the file server as the handler for
  // all URL paths that start with "/static/". For matching paths, we strip the
  // "/static" prefix before the request reaches the file server.
  mux.Get("/static/", http.StripPrefix("/static", fileServer))

  // Return the 'standard' middleware chain followed by the servemux.
  return standardMiddleware.Then(mux)
}
