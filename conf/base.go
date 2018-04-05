package conf
import(
	"encoding/json"
	"io/ioutil"
	"github.com/jmesyan/leaf/log"
)
var (
	LenStackBuf = 4096

	// log
	LogLevel string
	LogPath  string
	LogFlag  int
	LogMaxNumber int32
	LogMaxSize   int64
	LogUint      log.UNIT

	// console
	ConsolePort   int
	ConsolePrompt string = "Hewolf# "
	ProfilePath   string

	// cluster
	ListenAddr      string
	ConnAddrs       []string
	PendingWriteNum int
)

type BaseConf struct {
	LogLevel    string
	LogPath     string
	LogFlag    int
	LogMaxNumber int32
	LogMaxSize   int64
	LogUint      log.UNIT
	WSAddr      string
	CertFile    string
	KeyFile     string
	TCPAddr     string
	MaxConnNum  int
	ConsolePort int
	ProfilePath string
}

func globalConfInit(cconf *BaseConf){
	if cconf == nil {
		log.Fatal("the global conf  data is empty")
	}
	LogLevel = cconf.LogLevel
	LogPath = cconf.LogPath
	LogFlag = cconf.LogFlag
	LogMaxNumber = cconf.LogMaxNumber
	LogMaxSize = cconf.LogMaxSize
	LogUint = cconf.LogUint
	ConsolePort = cconf.ConsolePort
}

func NewBaseConf(path string) (*BaseConf, error){
	cconf := &BaseConf{}
	//集群服务器配置信息
	data, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, cconf)
	if err != nil {
		panic(err)
	}
	cfgpath = path

	globalConfInit(cconf)

	return cconf, nil
}
