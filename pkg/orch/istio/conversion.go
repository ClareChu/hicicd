package istio

import (
	"os"
	"path/filepath"
)

type Builder struct {
	Path       string
	Name       string
	FileType   string
	Profile    string
	ConfigType interface{}
}

func (b *Builder) WriterFile(by []byte)  (error){
	_, err := os.Stat(b.Path)
	if os.IsNotExist(err) {
		err = os.Mkdir(b.Path, os.ModePerm)
	}
	if err != nil {

	}
	fn := filepath.Join(b.Path, b.Name) + "." + b.FileType
	_, err = os.Stat(fn)
	if os.IsNotExist(err) {
		f, _ := os.OpenFile(fn, os.O_RDONLY | os.O_CREATE, 0666)
		f.Write(by)
		f.Close()
	}
	return nil
}

