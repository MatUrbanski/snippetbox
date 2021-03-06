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

  // Create a new middleware chain containing the middleware specific to
  // our dynamic application routes. For now, this chain will only contain
  // the session middleware but we'll add more to it later.
  dynamicMiddleware := alice.New(app.session.Enable, noSurf, app.authenticate)

  // Use the pat.New() function to initialize a new Pat router.
  mux := pat.New()

  // Register the home function as the handler for the "/" URL pattern.
  mux.Get("/", dynamicMiddleware.ThenFunc(app.home))

  // Register the createSnippetForm function as the handler for the GET "/snippet/create" URL pattern.
  mux.Get("/snippet/create", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.createSnippetForm))

  // Register the createSnippet function as the handler for the POST "/snippet/create" URL pattern.
  mux.Post("/snippet/create", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.createSnippet))

  // Register the showSnippet function as the handler for the "/snippet/:id" URL pattern.
  mux.Get("/snippet/:id", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.showSnippet))

  // Add routes for user signup, login and logout.
  mux.Get("/user/signup", dynamicMiddleware.ThenFunc(app.signupUserForm))
  mux.Post("/user/signup", dynamicMiddleware.ThenFunc(app.signupUser))
  mux.Get("/user/login", dynamicMiddleware.ThenFunc(app.loginUserForm))
  mux.Post("/user/login", dynamicMiddleware.ThenFunc(app.loginUser))
  mux.Post("/user/logout", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.logoutUser))

  // Add a new GET /ping route.
  mux.Get("/ping", http.HandlerFunc(ping))

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
