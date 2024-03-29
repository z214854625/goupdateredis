package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
)

var signalCh chan os.Signal
var IsProfile bool

func InstallSignal(callBack func()) {
	signalCh = make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGPIPE, syscall.SIGUSR1, syscall.SIGUSR2)
	procSignal(callBack)
}

func procSignal(callBack func()) {
	for {
		s := <-signalCh
		switch s {
		case syscall.SIGPIPE:
			continue
		case syscall.SIGTERM: //正常结束程序
			if callBack != nil {
				callBack()
			}
		case syscall.SIGUSR2: //启动profile分析
			if !IsProfile {
				go profile()
				fmt.PrintLn("profile() on 9913 success")
			}
		default:
			fmt.PrintLn("got sig %v", signalCh)
			if callBack != nil {
				callBack()
			}
		}
	}
}

func profile() { //开启profile
	IsProfile = true
	fmt.PrintLn(http.ListenAndServe("0.0.0.0:9913", nil))
}
