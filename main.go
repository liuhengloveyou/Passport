package main

import (
	"flag"
	"fmt"
	_ "net/http/pprof"
	"os"
	"runtime"
	"runtime/pprof"

	gocommon "github.com/liuhengloveyou/go-common"
	
	"github.com/liuhengloveyou/passport/common"
	"github.com/liuhengloveyou/passport/face"
)

var (
	BuildTime string
	CommitID  string

	showVer = flag.Bool("version", false, "show version")
	initSys = flag.Bool("init", false, "初始化系统.")

	cpuprofile = flag.String("cpuprofile", "", "write cpu profile `file`")
	memprofile = flag.String("memprofile", "", "write memory profile to `file`")
)

func main() {
	flag.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU())

	if *showVer {
		fmt.Printf("%s\t%s\n", BuildTime, CommitID)
		os.Exit(0)
	}

	gocommon.SingleInstane(common.ServConfig.PidFile)

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			panic(err)
		}
		if err := pprof.StartCPUProfile(f); err != nil {
			panic(err)
		}
		defer pprof.StopCPUProfile()
	}

	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {

		}
		runtime.GC() // get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			panic(err)
		}
		f.Close()
	}

	switch common.ServConfig.Face {
	case "http":
		face.InitAndRunHttpApi(&common.ServConfig)
	case "grpc":
		//face.GrpcFace()
	default:
		fmt.Println("face: [http | grpc]")
	}
}
