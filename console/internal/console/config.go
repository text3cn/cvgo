package console

import (
	"cvgo/console/internal/types"
	"cvgo/kit/filekit"
	"github.com/silenceper/log"
	"github.com/spf13/viper"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

var (
	HotCompileCfg   *types.HotCompileConfig
	CrossCompileCfg *types.CrossCompileCfg
	RootPath        string
	Exit            chan bool
	BuildPkg        string
	Started         chan bool
)

func LoadConfig() {
	HotCompileCfg = &types.HotCompileConfig{}
	CrossCompileCfg = &types.CrossCompileCfg{}
	RootPath = filekit.GetParentDir(3) + string(os.PathSeparator)

	// 默认配置
	setHotCompileDefaultConfig()
	setCrossCompileDefaultConfig()

	// 从配置文件加载配置
	filename := filepath.Join(RootPath, "console", "cvg.yaml")
	if filekit.FileExist(filename) {
		goodCfg := viper.New()
		goodCfg.AddConfigPath(filepath.Join(RootPath, "console"))
		goodCfg.SetConfigName("cvg")
		goodCfg.SetConfigType("yaml")
		if err := goodCfg.ReadInConfig(); err != nil {
			panic(err)
		}

		// 热编译
		if goodCfg.IsSet("hotCompilation.outputDir") {
			HotCompileCfg.HotCompileOutputDir = goodCfg.GetString("hotCompilation.outputDir")
		}
		if goodCfg.IsSet("hotCompilation.appName") {
			HotCompileCfg.AppName = goodCfg.GetString("hotCompilation.appName")
		}
		if goodCfg.IsSet("hotCompilation.watchExts") {
			HotCompileCfg.WatchExts = goodCfg.GetStringSlice("hotCompilation.watchExts")
		}
		if goodCfg.IsSet("hotCompilation.watchDirs") {
			HotCompileCfg.WatchPaths = goodCfg.GetStringSlice("hotCompilation.watchDirs")
		}
		if goodCfg.IsSet("hotCompilation.excludedDirs") {
			HotCompileCfg.ExcludedDirs = goodCfg.GetStringSlice("hotCompilation.excludedDirs")
		}
		if goodCfg.IsSet("hotCompilation.prevBuildCmds") {
			HotCompileCfg.ExcludedDirs = goodCfg.GetStringSlice("hotCompilation.prevBuildCmds")
		}

		// 交叉编译
		if goodCfg.IsSet("crossCompilation.outputDir") {
			CrossCompileCfg.OutputDir = goodCfg.GetString("crossCompilation.outputDir")
			if strings.HasPrefix(CrossCompileCfg.OutputDir, "./") ||
				(!strings.HasPrefix(CrossCompileCfg.OutputDir, "/") && !strings.HasPrefix(CrossCompileCfg.OutputDir, "./")) {
				CrossCompileCfg.OutputDir = filepath.Join(RootPath, CrossCompileCfg.OutputDir)
			}
		}
	}
}

// 热编译默认配置
func setHotCompileDefaultConfig() {
	if HotCompileCfg.AppName == "" {
		HotCompileCfg.AppName = path.Base(RootPath)
	}

	outputExt := ""
	if runtime.GOOS == "windows" {
		outputExt = ".exe"
	}
	if HotCompileCfg.HotCompileOutputDir == "" {
		HotCompileCfg.OutputAppPath = "./" + HotCompileCfg.AppName + outputExt
	} else {
		HotCompileCfg.OutputAppPath = HotCompileCfg.HotCompileOutputDir + string(filepath.Separator) + HotCompileCfg.AppName + outputExt
	}

	HotCompileCfg.WatchExts = append(HotCompileCfg.WatchExts, ".go")

	// set log level, default is debug
	if HotCompileCfg.LogLevel != "" {
		setLogLevel(HotCompileCfg.LogLevel)
	}
}

// 交叉编译默认配置
func setCrossCompileDefaultConfig() {
	if CrossCompileCfg.OutputDir == "" {
		CrossCompileCfg.OutputDir = filekit.GetParentDir(3) + string(os.PathSeparator) + "dist"
	}
	filekit.EnsureDirExists(CrossCompileCfg.OutputDir)
}

func setLogLevel(level string) {
	switch level {
	case "debug":
		log.SetLogLevel(log.LevelDebug)
	case "info":
		log.SetLogLevel(log.LevelInfo)
	case "warn":
		log.SetLogLevel(log.LevelWarning)
	case "error":
		log.SetLogLevel(log.LevelError)
	case "fatal":
		log.SetLogLevel(log.LevelFatal)
	default:
		log.SetLogLevel(log.LevelDebug)
	}
}
