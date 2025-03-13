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
	InfoLog  *log.Logger
	ErrorLog *log.Logger
	RootPath string
	Render   *render.Render
	Routes   *chi.Mux
	config   config
}

type config struct {
	port     string
	renderer string
}

func (q *Quokka) New(rootPath string) error {
	pathConfigs := initPaths{
		rootPath:    rootPath,
		folderNames: []string{"handlers", "migrations", "data", "views", "public", "tmp", "logs", "middlewares"},
	}
	err := q.Init(pathConfigs)
	if err != nil {
		return err
	}

	// check .env file exist
	err = q.checkDotEnvExist(rootPath)
	if err != nil {
		return err
	}

	// load .env file
	err = godotenv.Load(rootPath + "/.env")
	if err != nil {
		return err
	}

	// create loggers and another configuration of quokka framework
	infoLog, errorLog := q.startLoggers()
	q.InfoLog = infoLog
	q.ErrorLog = errorLog
	q.Debug, _ = strconv.ParseBool(os.Getenv("DEBUG"))
	q.Version = version
	q.RootPath = rootPath
	q.Routes = q.routes().(*chi.Mux)

	q.config = config{
		port:     os.Getenv("PORT"),
		renderer: os.Getenv("RENDERER"),
	}

	q.Render = q.createRenderer(q)

	return nil
}

func (q *Quokka) Init(p initPaths) error {
	rootPath := p.rootPath
	for _, path := range p.folderNames {
		err := q.createDirIfNotExist(rootPath + "/" + path)
		if err != nil {
			return err
		}
	}
	return nil
}

func (q *Quokka) checkDotEnvExist(path string) error {
	err := q.createFileIfNotExists(fmt.Sprintf("%s/.env", path))
	if err != nil {
		return err
	}
	return nil
}

func (q *Quokka) startLoggers() (*log.Logger, *log.Logger) {
	var infoLog *log.Logger
	var errorLog *log.Logger

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	return infoLog, errorLog
}

// listening to server
func (q *Quokka) ListenAndServe() {
	serv := &http.Server{
		Addr:         fmt.Sprintf(":%s", os.Getenv("PORT")),
		Handler:      q.Routes,
		IdleTimeout:  30 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 600 * time.Second,
	}
	q.InfoLog.Printf("Listening to port: %s", os.Getenv("PORT"))
	err := serv.ListenAndServe()
	q.ErrorLog.Fatal(err)
}

// create render engine
func (q *Quokka) createRenderer(quo *Quokka) *render.Render {
	myRenderer := render.Render{
		RootPath: quo.RootPath,
		Port:     quo.config.port,
		Renderer: quo.config.renderer,
	}
	return &myRenderer
}
