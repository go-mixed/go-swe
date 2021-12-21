package main

import (
	"github.com/spf13/cobra"
	web2 "go-common-web"
	"go-common/task_pool"
	"go-common/utils/io"
	"go-swe/src/web"

	"fmt"
	"go-common/utils"
	"go-swe/src/astro"
	conf "go-swe/src/settings"
	"path/filepath"
	"time"
)

var settings *conf.Settings

func main() {
	// 读取当前执行文件的目录
	currentDir := io_utils.GetCurrentDir()

	rootCmd := &cobra.Command{
		Use:   "go-swe-server",
		Short: "a server of Swiss Ephemeris",
		Run: func(cmd *cobra.Command, args []string) {
			config, _ := cmd.PersistentFlags().GetString("config")
			log, _ := cmd.PersistentFlags().GetString("log")
			run(config, log)
		},
	}

	// 读取CLI
	rootCmd.PersistentFlags().StringP("config", "c", filepath.Join(currentDir, "conf/settings.json"), "config file")
	rootCmd.PersistentFlags().StringP("log", "l", filepath.Join(currentDir, "logs"), "log files path")
	err := rootCmd.Execute()
	if err != nil {
		panic(err.Error())
	}
}

func run(_configFile, _logPath string) {
	var err error
	// 读取配置文件
	settings, err = conf.LoadSettings(_configFile)
	if err != nil {
		panic(err.Error())
	}
	if settings == nil {
		panic("read settings fatal.")
	}

	// 初始化日志
	utils.InitLogger(filepath.Join(_logPath, "app.log"), "")
	logger := utils.GetSugaredLogger()

	exec, _ := task_pool.NewExecutor(task_pool.DefaultExecutorParams(), logger)
	defer exec.Stop()
	exec.ListenStopSignal()

	exec.Submit(func(stopChan <-chan struct{}) {
		engine := web2.NewGinEngine(web2.DefaultGinOptions(settings.Debug, true))
		web.RegisterRouter(engine)
		server := web2.NewHttpServer(web2.DefaultServerOptions(settings.Host))
		if err := server.SetDefaultServeHandler(engine, web2.NewCertificate(settings.Cert, settings.Key)); err != nil {
			logger.Error(err.Error())
			return
		}
		if err := server.Run(stopChan, nil); err != nil {
			logger.Error(err.Error())
		}
	})

	exec.Wait()
	logger.Info("main application exit.")

	t := time.Now().UnixNano()

	long, _ := astro.StringToDegrees("116°23'")
	lat, _ := astro.StringToDegrees("39°54'")
	fmt.Printf("Geo: %f %f\n", long, lat)

	geo := &astro.GeographicCoordinates{
		Longitude: astro.ToRadians(long),
		Latitude:  astro.ToRadians(lat),
	}
	tz, _ := time.LoadLocation("Asia/Shanghai")

	year, month, day := time.Now().Date()

	astronomy := astro.NewAstronomy()

	t = time.Now().UnixNano()

	// 东八区的正午是UTC的4点
	noon := time.Date(year, month, day, 4, 0, 0, 0, time.UTC)
	jd := astro.TimeToJulianDay(noon)
	deltaT := astro.DeltaT(jd)
	et := jd.Add(deltaT)
	etT := et.ToTime(time.UTC)
	fmt.Printf("JD: %f at %v \n", jd, jd.ToTime(time.UTC))
	fmt.Printf("ET: %f at %v deltaT: %v\n", et, etT, deltaT)

	// 太阳
	sunTimes, err := astronomy.SunTwilight(jd, geo, false)
	if err != nil {
		fmt.Printf("SunTwilight Error: %s", err.Error())
	}
	fmt.Printf("Sun Rise: %v\n", sunTimes.Rise.ToTime(tz))
	fmt.Printf("Sun Set: %v\n", sunTimes.Set.ToTime(tz))
	fmt.Printf("Sun Culmination: %v | %v\n", sunTimes.Culmination.ToTime(tz), sunTimes.LowerCulmination.ToTime(tz))
	fmt.Printf("Sun Civil : %v | %v\n", sunTimes.Civil.Dawn.ToTime(tz), sunTimes.Civil.Dusk.ToTime(tz))
	fmt.Printf("Night : %v\n", sunTimes.Night())

	fmt.Printf("SunTwilight: %.4f ms\n----------------------\n", float64(time.Now().UnixNano()-t)/1e6)
	t = time.Now().UnixNano()

	// 月亮
	moonTimes, err := astronomy.MoonTwilight(jd, geo, false)
	if err != nil {
		fmt.Printf("MoonTwilight Error: %s", err.Error())
	}
	fmt.Printf("Moon Rise: %v\n", moonTimes.Rise.ToTime(tz))
	fmt.Printf("Moon Set: %v\n", moonTimes.Set.ToTime(tz))
	fmt.Printf("Moon Culmination: %v | %v\n", moonTimes.Culmination.ToTime(tz), moonTimes.LowerCulmination.ToTime(tz))

	fmt.Printf("MoonTwilight: %.4f ms\n----------------------\n", float64(time.Now().UnixNano()-t)/1e6)
	t = time.Now().UnixNano()

}
