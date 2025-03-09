package quokka

import (
	"fmt"
	"log"
	"os"
	"strconv"

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

	// create loggers
	infoLog, errorLog := q.startLoggers()
	q.InfoLog = infoLog
	q.ErrorLog = errorLog
	q.Debug, _ = strconv.ParseBool(os.Getenv("DEBUG"))
	q.Version = version

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
