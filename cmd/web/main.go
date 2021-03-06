package main

import (
  "crypto/tls"
  "database/sql"
  "flag"
  "html/template"
  "log"
  "net/http"
  "os"
  "time"

  "mateuszurbanski/snippetbox/pkg/models"
  "mateuszurbanski/snippetbox/pkg/models/mysql"

  _ "github.com/go-sql-driver/mysql"
  "github.com/golangcollege/sessions"
)

// Notice how the import path for our driver is prefixed with an underscore? This is because
// our main.go file doesn’t actually use anything in the mysql package. So if we try to
// import it normally the Go compiler will raise an error. However, we need the driver’s
// init() function to run so that it can register itself with the database/sql package. The
// trick to getting around this is to alias the package name to the blank identifier. This is
// standard practice for most of Go’s SQL drivers.

type contextKey string

const contextKeyIsAuthenticated = contextKey("isAuthenticated")

// Define an application struct to hold the application-wide dependencies.
type application struct {
  errorLog *log.Logger
  infoLog  *log.Logger
  session  *sessions.Session
  snippets interface {
    Insert(string, string, string) (int, error)
    Get(int) (*models.Snippet, error)
    Latest() ([]*models.Snippet, error)
  }
  templateCache map[string]*template.Template
  users         interface {
    Insert(string, string, string) error
    Authenticate(string, string) (int, error)
    Get(int) (*models.User, error)
  }
}

func main() {
  // Define a new command-line flag with the name 'addr', a default value of ":4000"
  // and some short help text explaining what the flag controls. The value of the
  // flag will be stored in the addr variable at runtime.
  addr := flag.String("addr", ":4000", "HTTP network address")

  // Define a new command-line flag for the MySQL DSN string.
  dsn := flag.String("dsn", "web:password@/snippetbox?parseTime=true", "MySQL data source name")

  // Define a new command-line flag for the session secret (a random key which
  // will be used to encrypt and authenticate session cookies). It should be 32
  // bytes long.
  secret := flag.String("secret", "n6Gdh+pPbnzHbS*+9Pk8qGWhTzbpa@gd", "Secret key")

  // Define a new command-line flag for the current environment.
  environment := flag.String("environment", "development", "Current Environment")


  // Importantly, we use the flag.Parse() function to parse the command-line flag.
  // This reads in the command-line flag value and assigns it to the addr
  // variable. You need to call this *before* you use the addr variable
  // otherwise it will always contain the default value of ":4000". If any errors are
  // encountered during parsing the application will be terminated.
  flag.Parse()

  // Use log.New() to create a logger for writing information messages. This takes
  // three parameters: the destination to write the logs to (os.Stdout), a string
  // prefix for message (INFO followed by a tab), and flags to indicate what
  // additional information to include (local date and time). Note that the flags
  // are joined using the bitwise OR operator |.
  infoLog := log.New(os.Stdout, "Info\t", log.Ldate|log.Ltime)

  // Create a logger for writing error messages in the same way, but use stderr as
  // the destination and use the log.Lshortfile flag to include the relevant
  // file name and line number.
  errorLog := log.New(os.Stderr, "Error\t", log.Ldate|log.Ltime|log.Lshortfile)

  // To keep the main() function tidy I've put the code for creating a connection
  // pool into the separate openDB() function below. We pass openDB() the DSN
  // from the command-line flag.
  db, err := openDB(*dsn)
  if err != nil {
    errorLog.Fatal(err)
  }

  // We also defer a call to db.Close(), so that the connection pool is closed
  // before the main() function exits.
  defer db.Close()

  // Initialize a new template cache...
  templateCache, err := newTemplateCache("./ui/html/")

  if err != nil {
    errorLog.Fatal(err)
  }

  // Use the sessions.New() function to initialize a new session manager,
  // passing in the secret key as the parameter. Then we configure it so
  // sessions always expires after 12 hours.
  session := sessions.New([]byte(*secret))
  session.Lifetime = 12 * time.Hour
  session.Secure = true // Set the Secure flag on our session cookies

  // Initialize a new instance of application containing the dependencies.
  app := &application{
    errorLog:      errorLog,
    infoLog:       infoLog,
    session:       session,
    snippets:      &mysql.SnippetModel{DB: db},
    templateCache: templateCache,
    users:         &mysql.UserModel{DB: db},
  }

  // Initialize a tls.Config struct to hold the non-default TLS settings we want
  // the server to use.
  tlsConfig := &tls.Config{
    PreferServerCipherSuites: true,
    CurvePreferences:         []tls.CurveID{tls.X25519, tls.CurveP256},
  }

  // Initialize a new http.Server struct. We set the Addr and Handler fields so
  // that the server uses the same network address and routes as before, and set
  // the ErrorLog field so that the server now uses the custom errorLog logger in
  // the event of any problems.
  srv := &http.Server{
    Addr:         *addr,
    ErrorLog:     errorLog,
    Handler:      app.routes(),
    TLSConfig:    tlsConfig,
    IdleTimeout:  time.Minute,
    ReadTimeout:  5 * time.Second,
    WriteTimeout: 10 * time.Second,
  }

  // Use the http.ListenAndServe() function to start a new web server. We pass in
  // two parameters: the TCP network addres to listen on (in this case ":4000")
  // and the servemux we just created. If http.ListenAndServe() returns an error
  // we use log.Fatal() function to log the error message and exit. Note
  // that any error returned by http.ListenAndServe() is always non-nil.
  // The value returned from the flag.String() function is a pointer to the flag
  // value, not the value itself. So we need to dereference the pointer (i.e.
  // prefix it with the * symbol) before using it. Note that we're using the
  // log.Printf() function to interpolate the address with the log message.
  infoLog.Printf("Starting server on %s", *addr)


  if *environment == "production" {
    err = srv.ListenAndServe()
  } else {
    err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
  }

  errorLog.Fatal(err)
}

// The openDB() function wraps sql.Open() and returns a sql.DB connection pool
// for a given DSN.
func openDB(dsn string) (*sql.DB, error) {
  db, err := sql.Open("mysql", dsn)
  if err != nil {
    return nil, err
  }

  if err = db.Ping(); err != nil {
    return nil, err
  }

  return db, nil
}
