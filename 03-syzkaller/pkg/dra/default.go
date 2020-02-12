package dra

import (
	"os"
	"path/filepath"
)

// useful const
const (
	//startTime = 10800
	startTime       = 0
	newTime         = 3600
	bootTime        = 3600
	TimeWriteToDisk = 3600
	TimeExit        = 3600 * 24

	TaskNum             = 40
	TaskCountLimitation = 20
	TaskBase            = 1

	DebugLevel = 2

	CollectPath     = true
	CollectUnstable = true

	// collect coverage by intersection instead of union.
	StableCoverage = true
	// check Condition address in syz-fuzzer
	CheckCondition = true
)

var pathHome = os.Getenv("HOME")
var pathRoot = filepath.Join(pathHome, "data")

var pathLinux = filepath.Join(pathRoot, "benchmark/linux/13-linux-clang-np")
var FileVmlinuxObjdump = filepath.Join(pathLinux, "vmlinux.objdump")

var pathGit = filepath.Join(pathRoot, "git")
var pathRepo = filepath.Join(pathGit, "gopath/src/github.com/ZHYfeng/2018_dependency")
var PathA2i = filepath.Join(pathRepo, "02-dependency/cmake-build-debug/tools/A2I/a2i")

const (
	NameDevice         = "dev_"
	NameBase           = "base"
	NameWithDra        = "01-result-with-dra"
	NameWithoutDra     = "02-result-without-dra"
	NameData           = "data.txt"
	NameDataDependency = "dataDependency.bin"
	NameDataResult     = "dataResult.bin"
	NameDataRunTime    = "dataRunTime.bin"
	NameStatistics     = "statistics.bin"
	NameUnstable       = "unstable.bin"
	NameUnstableResult = "unstable.txt"

	NameDriver   = "built-in"
	FileAsm      = NameDriver + ".s"
	FileTaint    = NameDriver + ".taint"
	FileFunction = NameDriver + ".function.json"
	FileBc       = NameDriver + ".bc"
)
