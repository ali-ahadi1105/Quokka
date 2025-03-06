package quokka

const version = "1.0.0"

type Quokka struct {
	AppName string
	Debug   bool
	Version string
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
