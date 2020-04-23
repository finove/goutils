package main

import (
	"flag"
	"github.com/finove/goutils/logger"
	"github.com/finove/goutils/sms"
	"github.com/finove/goutils/vconfig"
)

var goutilsVersion = "1.0.0"

var testFor string

var (
	smsKeyID             string
	smsSecret            string
	smsSignName          string
	smsVerifyTemplate    string
	smsCode, smsPhoneNum string
)

var (
	configFileName string
)

func main() {
	parseCmdline()
	switch testFor {
	case "log":
		testLog()
	case "cfg":
		testCfg()
	default:
		testLog()
	}
}

func parseCmdline() {
	flag.StringVar(&testFor, "t", "", "what you want test")
	flag.StringVar(&configFileName, "cfg", "", "config file to load")
	flag.Parse()
}

func testLog() {
	logger.Setup(true, "debug", "")
	logger.SetLevel("info")
	logger.Fatal("this is fatal log %d", 1)
	logger.Error("this is error log %d", 2)
	logger.Warning("this is warning log %d", 3)
	logger.Info("this is info log %d", 4)
	logger.Debug("this is debug log %d", 5)
}

func testSms() {
	var err error
	sms.AliService.ConfigAuth("", smsKeyID, smsSecret)
	sms.AliService.SetupVerifyCode(smsSignName, smsVerifyTemplate)
	_, err = sms.AliService.SendVerifyCode(smsCode, smsPhoneNum)
	logger.Info("send verify code %s", err)
}

func testCfg() {
	logger.Info("load config1 %v", vconfig.LoadConfigFile(configFileName, true))
	vconfig.Viper.SetDefault("cfg1", "value1")
	logger.Info("get cfg1 = %s", vconfig.Viper.GetString("cfg1"))
	vconfig.Viper.WriteConfig()
}
