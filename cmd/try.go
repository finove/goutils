package cmd

import (
	"fmt"

	"github.com/finove/goutils/db/redisop"
	"github.com/finove/goutils/errormessage"
	"github.com/finove/goutils/logger"
	"github.com/finove/goutils/vconfig"
	"github.com/garyburd/redigo/redis"
	"github.com/spf13/cobra"
)

var (
	tryWhat        string
	tryCfgFileName string
)

var tryWhatList = []string{"log", "cfg", "errmsg", "redis"}

var tryCmd = &cobra.Command{
	Use:   "try",
	Short: "try functions",
	Long:  `try test functions`,
	Run: func(cmd *cobra.Command, args []string) {
		switch tryWhat {
		case "log":
			testLog()
		case "cfg":
			testCfg()
		case "errmsg":
			testErrorMsg()
		case "redis":
			testRedis()
		default:
			testLog()
		}
	},
}

func init() {
	rootCmd.AddCommand(tryCmd)
	tryCmd.Flags().StringVar(&tryWhat, "what", "", "what you want test")
	tryCmd.Flags().StringVar(&tryCfgFileName, "cfg", "", "config file to load for test")
	tryCmd.MarkFlagRequired("what")
}

func testLog() {
	logger.Setup(true, "debug", "t.log", `{"tosyslog":true, "appname":"goutils"}`)
	logger.SetLevel("info")
	logger.Fatal("this is fatal log %d", 1)
	logger.Error("this is error log %d", 2)
	logger.Warning("this is warning log %d", 3)
	logger.Info("this is info log %d", 4)
	logger.Debug("this is debug log %d", 5)
	logger.Info("support trywhat %v", tryWhatList)
}

func testErrorMsg() {
	errormessage.AddErrorMessages(map[int][]string{
		2001: {"403", "fail3"},
		2002: {"404", "fail4"},
		2003: {"494", "fail5"},
	})
	errormessage.AddErrorMessage(1002, "401", "未登录", "认证失败或没有授权")
	logger.Info("code 0 = %d %s", errormessage.HTTPStatus(0), errormessage.Message(0))
	logger.Info("code 1001 = %d %s", errormessage.HTTPStatus(1001), errormessage.Message(1001, fmt.Errorf("more %s", "info")))
	logger.Info("code 1002 = %d %s", errormessage.HTTPStatus(1002), errormessage.Message(1002, fmt.Errorf("补充失败信息 %d", 1002)))
	logger.Info("code 2001 = %d %s", errormessage.HTTPStatus(2001), errormessage.Message(2001))
	logger.Info("code 2002 = %d %s", errormessage.HTTPStatus(2002), errormessage.Message(2002))
	logger.Info("code 2003 = %d %s", errormessage.HTTPStatus(2003), errormessage.Message(2003))
}

func testCfg() {
	logger.Info("load config1 %v", vconfig.LoadConfigFile(tryCfgFileName, true))
	vconfig.Viper.SetDefault("cfg1", "value1")
	logger.Info("get cfg1 = %s", vconfig.Viper.GetString("cfg1"))
	vconfig.Viper.WriteConfig()
}

func testRedis() {
	var r = &redisop.RedisPoolX{}
	testLog()
	logger.Info("test redisop")
	redisop.SetLogger(logger.GetLogger())
	r.Connect("default1")
	r.Rdo("SET", "tmpa", "HELLO WORLD")
	vv, err := redis.String(r.Rdo("GET", "tmpa"))
	fmt.Printf("err %v, value %v\n", err, vv)
	fmt.Printf("redis notify config %s\n", r.GetNotifyConfig())
}
