package main

import (
  "errors"
  "fmt"
  "net/http"
  "strconv"
  "mateuszurbanski/snippetbox/pkg/models"
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

  s, err := app.snippets.Latest()
  if err != nil {
    app.serverError(w, err)
    return
  }

  // Use the new render helper.
  app.render(w, r, "home.page.tmpl", &templateData{Snippets: s,})
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

  // Use the SnippetModel object's Get method to retrieve the data for a
  // specific record based on its ID. If no matching record is found,
  // return a 404 Not Found response.
  s, err := app.snippets.Get(id)

  if err != nil {
    if errors.Is(err, models.ErrNoRecord) {
      app.notFound(w)
    } else {
      app.serverError(w, err)
    }

    return
  }

  app.render(w, r, "show.page.tmpl", &templateData{Snippet: s})
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

  // Create some variables holding dummy data. We'll remove these later on
  // during the build.
  title := "O snail"
  content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
  expires := "7"

  // Pass the data to the SnippetModel.Insert() method, receiving the
  // ID of the new record back.
  id, err := app.snippets.Insert(title, content, expires)

  if err != nil {
    app.serverError(w, err)
    return
  }

  // Redirect the user to the relevant page for the snippet.
  http.Redirect(w, r, fmt.Sprintf("/snippet?id=%d", id), http.StatusSeeOther)
}
