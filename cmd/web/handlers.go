package main

import (
  "fmt"
  "html/template"
  "net/http"
  "strconv"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
  // Check if the current request URL path exactly matches "/". If it doesn't, use
  // the http.NotFound() function to send a 404 response to the client.
  // Importantly, we then return from the handler. If we don't return, the handler
  // would keep execurint and also write the "Hello from Snippetbox".
  if r.URL.Path != "/" {
    app.notFound(w) // Use the notFound() helper
    return
  }

  // Initialize a slice containing the paths to the two files. Note that the
  // home.page.tmpl file must be the *first* file in the slice.
  files := []string{
    "./ui/html/home.page.tmpl",
    "./ui/html/base.layout.tmpl",
    "./ui/html/footer.partial.tmpl",
  }

  // Use the template.ParseFiles() function to read the template file into a
  // template set. If there's an error, we log the detailed error message and use
  // the app.serverError() function to send a generic 500 Internal Server Error
  // response to the user.
  ts, err := template.ParseFiles(files...)
  if err != nil {
    app.serverError(w, err) // Use the serverError() helper.
    return
  }

  // We then use the Execute() method on the template set to write the template
  // content as the response body. The last parameter to Execute() represents any
  // dynamic data that we want to pass in, which for now we'll leave as nil.
  err = ts.Execute(w, nil)
  if err != nil {
    // Because the home handler function is now a method against application
    // it can access its fields, including the error logger. We'll write the log
    // message to this instead of the standard logger.
    app.serverError(w, err) // Use the serverError() helper.
  }
}

func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
  // Extract the value of the id parameter from the query string and try to
  // convert it to an integer using the strconv.Atoi(i) function. If it can't
  // be converted to an integer, or the value is less than 1, we return a 404 page
  // not found response.
  id, err := strconv.Atoi(r.URL.Query().Get("id"))

  if err != nil || id < 1 {
    app.notFound(w) // Use the notFound() helper.
    return
  }

  // Use the fmt.Fprintf() function to interpolate the id value with our response
  // and write it to the http.ResponseWriter.
  fmt.Fprintf(w, "Display a specific snippet with ID %d", id)
}

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
  // Use r.Method to check whether the request is using POST or not. Note that
  // http.MethodPost is a constant equal to the string "POST".
  if r.Method != http.MethodPost {
    // Use the Header().Set() method to add an 'Allow: POST' header to the
    // response header map. The first paramter is the header name, and
    // the second parameter is the header value.
    w.Header().Set("Allow", http.MethodPost)

    // Use the http.Error() function to send a 405 status code and "Method not
    // Allowed" string as the response body.
    app.clientError(w, http.StatusMethodNotAllowed) // Use the clientError() helper.
    return
  }

  w.Write([]byte("Create a specific snippet..."))
}
