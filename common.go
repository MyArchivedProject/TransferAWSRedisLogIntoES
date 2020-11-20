package main

import (
	"log"
	"os"
	"runtime"
)

func errorExit(err error) {
	if err != nil {
		pc, _, _, _ := runtime.Caller(2)
		funcName := runtime.FuncForPC(pc).Name() // 获取函数调用者的名字
		log.Println("Error, exit.\n Function: " + funcName + "Error:\n" + error.Error(err))
		os.Exit(1)
	}
}

func errorTolerate(err error) {
	if err != nil {
		pc, _, _, _ := runtime.Caller(2)
		funcName := runtime.FuncForPC(pc).Name() // 获取函数调用者的名字
		log.Println("Error, but tolerate it.\n Function: " + funcName + "Error:\n" + error.Error(err))
	}
}
func printLog(smg interface{}) {
	pc, _, _, _ := runtime.Caller(2)
	funcName := runtime.FuncForPC(pc).Name() // 获取函数调用者的名字
	log.Println("Function: " + funcName + "  message: ")
	log.Println(smg)
}
