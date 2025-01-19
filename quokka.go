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
	ErrorLog *log.Logger
	InfoLog  *log.Logger
	RootPath string
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

	// create logs
	infoLog, errorLog := quokka.startLoggers()
	quokka.InfoLog = infoLog
	quokka.ErrorLog = errorLog
	quokka.Debug, _ = strconv.ParseBool(os.Getenv("DEBUG"))
	quokka.Version = version
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
