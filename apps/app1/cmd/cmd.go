package cmd

import "flag"

var ENV string
var SIG string
var DAEMON string

//命令行处理参数
func Init() {
	flag.StringVar(&ENV, "env", "debug", "env")
	flag.StringVar(&SIG, "s", "start", "start:启动 reload:重载 stop:停止")
	flag.StringVar(&DAEMON, "daemon", "false", "true:daemon模式运行 false:前台模式运行")
	flag.Parse()
}
