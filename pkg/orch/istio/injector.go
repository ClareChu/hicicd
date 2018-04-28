package istio

import (
	"time"
	"istio.io/istio/pilot/pkg/model"
	"istio.io/istio/pilot/pkg/kube/inject"
	"os"
	"gopkg.in/yaml.v2"
	"bytes"
<<<<<<< HEAD
	"path/filepath"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/hiboot/pkg/utils"
)

type resource struct{
		enableAuth          bool
		in                  string
		want                string
		imagePullPolicy     string
		enableCoreDump      bool
		debugMode           bool
		duration            time.Duration
		includeIPRanges     string
		excludeIPRanges     string
		includeInboundPorts string
		excludeInboundPorts string
		tproxy              bool
=======
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/hiboot/pkg/utils"
	"path/filepath"
)

type resource struct {
	enableAuth          bool
	in                  string
	want                string
	imagePullPolicy     string
	enableCoreDump      bool
	debugMode           bool
	duration            time.Duration
	includeIPRanges     string
	excludeIPRanges     string
	includeInboundPorts string
	excludeInboundPorts string
	tproxy              bool
>>>>>>> master
}

const (
	DefaultSidecarProxyUID     = uint64(1337)
	DefaultVerbosity           = 2
	DefaultImagePullPolicy     = "IfNotPresent"
	DefaultIncludeIPRanges     = "*"
	DefaultIncludeInboundPorts = "*"
)

type Injector struct {
	Path       string
	Name       string
	FileType   string
	Profile    string
	Version    string
	ConfigType interface{}
}

const (
	Hub = "docker.io/istio"

	Tag = "latest"
)

<<<<<<< HEAD

func (i *Injector) IntoResource( cfg interface{}) (error){
=======
func (i *Injector) IntoResource(cfg interface{}) (error) {
>>>>>>> master
	//var dc v1.DeploymentConfig
	debugMode := true
	mesh := model.DefaultMeshConfig()
	c := resource{
<<<<<<< HEAD
		in:                  i.Path,
=======
		in:                  filepath.Join(i.Path, i.Name),
		want:                filepath.Join(i.Path, i.Name + "dc"),
>>>>>>> master
		includeIPRanges:     DefaultIncludeIPRanges,
		includeInboundPorts: DefaultIncludeInboundPorts,
	}

	params := &inject.Params{
<<<<<<< HEAD
		InitImage:           inject.InitImageName(Hub, Tag, debugMode),
		ProxyImage:          inject.ProxyImageName(Hub, Tag, debugMode),
		ImagePullPolicy:     "IfNotPresent",
		Verbosity:           DefaultVerbosity,
		SidecarProxyUID:     DefaultSidecarProxyUID,
		Version:             i.Version,
		EnableCoreDump:      false,
		Mesh:                &mesh,
		DebugMode:           debugMode,
	}
	sidecarTemplate, err := inject.GenerateTemplateFromParams(params)
	if err !=nil {
		return  err
	}
	in, err := os.Open(c.in)
	if err != nil {
		return  err
	}
	defer func() { _ = in.Close() }()
	var  out bytes.Buffer
	if err = inject.IntoResourceFile(sidecarTemplate, &mesh, in, &out); err != nil {
		return  err
	}
	in1 := out.Bytes()
	err = yaml.Unmarshal(in1, &cfg)
	if err !=nil {
		return  err
	}
	//privileged := true/
	//cfg.Spec.Template.Spec.InitContainers[0].SecurityContext.Privileged = &privileged
	return  nil
}


func InjectSideCar(cfg interface{}, fullName, version string)  (error){
	in, err :=	yaml.Marshal(cfg)
	if err != nil {
		log.Print("test to yml",err)
		return err
	}
utils.
	err = utils.WriterFile(in)
	if err != nil {
		log.Print("test to yml",err)
		return err
	}
	path := filepath.Join(os.TempDir(), fullName + ".yml")
	err = IntoResource(path, version, cfg)
	if err != nil {
		log.Print("test to yml",err)
	}
	return nil
}




=======
		InitImage:       inject.InitImageName(Hub, Tag, debugMode),
		ProxyImage:      inject.ProxyImageName(Hub, Tag, debugMode),
		ImagePullPolicy: "IfNotPresent",
		Verbosity:       DefaultVerbosity,
		SidecarProxyUID: DefaultSidecarProxyUID,
		Version:         i.Version,
		EnableCoreDump:  false,
		Mesh:            &mesh,
		DebugMode:       debugMode,
	}
	sidecarTemplate, err := inject.GenerateTemplateFromParams(params)
	if err != nil {
		return err
	}
	in, err := os.Open(c.in)
	if err != nil {
		return err
	}
	defer func() { _ = in.Close() }()
	var out bytes.Buffer
	if err = inject.IntoResourceFile(sidecarTemplate, &mesh, in, &out); err != nil {
		return err
	}
	in1 := out.Bytes()
	err = yaml.Unmarshal(in1, &cfg)
	log.Infof("",cfg)
	if err != nil {
		return err
	}
	return nil
}

func InjectSideCar(cfg interface{}, fullName, version string) (error) {
	in, err := yaml.Marshal(cfg)
	inje := Injector{
		Path:    os.TempDir(),
		Name:    fullName,
		Version: version,
	}
	if err != nil {
		log.Print("test to yml", err)
		return err
	}
	t, err := utils.WriterFile(inje.Path, inje.Name, in)
	log.Println("输出流", t)
	if err != nil {
		log.Print("test to yml", err)
		return err
	}
	err = inje.IntoResource(cfg)
	if err != nil {
		log.Print("test to yml", err)
	}
	return nil
}
>>>>>>> master
