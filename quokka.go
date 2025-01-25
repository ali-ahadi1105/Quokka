package quokka

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/ali-ahadi1105/Quokka/render"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

const version = "1.0.0"

type Quokka struct {
	AppName  string
	Debug    bool
	Version  string
	ErrorLog *log.Logger
	InfoLog  *log.Logger
	RootPath string
	Routes   *chi.Mux
	config   config
	Render   *render.Renderer
}

type config struct {
	port     string
	renderer string
}

func (quokka *Quokka) New(rootPath string) error {
	initConfigs := initPaths{
		rootPath:    rootPath,
		folderNames: []string{"handlers", "middlewares", "logs", "public", "data", "migrations", "views", "tmp"},
	}

	// create folders structur of application
	err := quokka.Init(initConfigs)

	if err != nil {
		return err
	}

	// check .env exist
	err = quokka.checkDotenvFile(rootPath)
	if err != nil {
		return err
	}

	// load env file
	err = godotenv.Load(rootPath + "/.env")
	if err != nil {
		return err
	}

	// create logs and other things
	infoLog, errorLog := quokka.startLoggers()
	quokka.InfoLog = infoLog
	quokka.ErrorLog = errorLog
	quokka.Debug, _ = strconv.ParseBool(os.Getenv("DEBUG"))
	quokka.Version = version
	quokka.RootPath = rootPath
	quokka.Routes = quokka.routes().(*chi.Mux)

	// create configuration files
	quokka.config = config{
		port:     os.Getenv("PORT"),
		renderer: os.Getenv("RENDERER"),
	}

	quokka.Render = quokka.createRenderer(quokka)

	return nil
}

func (quokka *Quokka) Init(path initPaths) error {
	root := path.rootPath
	for _, path := range path.folderNames {
		err := quokka.createDirIfNotExist(root + "/" + path)
		if err != nil {
			return err
		}
	}
	return nil
}

// start server
func (quo *Quokka) ListenAndServe() {
	server := http.Server{
		Addr:         fmt.Sprintf(":%s", os.Getenv("PORT")),
		Handler:      quo.Routes,
		ErrorLog:     quo.ErrorLog,
		IdleTimeout:  30 * time.Second,
		WriteTimeout: 600 * time.Second,
		ReadTimeout:  30 * time.Second,
	}

	quo.InfoLog.Printf("Server is running on port %s", os.Getenv("PORT"))
	err := server.ListenAndServe()
	quo.ErrorLog.Fatal(err)
}

func (quokka *Quokka) checkDotenvFile(path string) error {
	err := quokka.createFileIfNotExists(fmt.Sprintf("%s/.env", path))
	if err != nil {
		return err
	}
	return nil
}

func (quokka *Quokka) startLoggers() (*log.Logger, *log.Logger) {
	var infoLog *log.Logger
	var errorLog *log.Logger

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	return infoLog, errorLog
}

func (q *Quokka) createRenderer(quo *Quokka) *render.Renderer {
	myRenderer := render.Renderer{
		Renderer: quo.config.renderer,
		Port:     quo.config.port,
		RootPath: quo.RootPath,
	}
	return &myRenderer
}
