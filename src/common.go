package src

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"time"
)
// PrintErrorExit 我的不可原谅的错误
func PrintErrorExit(err error) {
	if err != nil {
		pc, _, _, _ := runtime.Caller(1)
		funcName := runtime.FuncForPC(pc).Name() // 获取函数调用者的名字
		log.Println("-----------------------------\nError, exit.\n Function: " + funcName + "Error:\n" + error.Error(err))
		os.Exit(1)
	}
}
// PrintErrorTolerate 我的可原谅的错误
func PrintErrorTolerate(err error) {
	if err != nil {
		pc, _, _, _ := runtime.Caller(1)
		funcName := runtime.FuncForPC(pc).Name() // 获取函数调用者的名字
		log.Println("Error, but tolerate it.\n Function: " + funcName + "Error:\n" + error.Error(err))
	}
}
// PrintLog 我的日志输出方式
func PrintLog(smg interface{}) {
	pc, _, _, _ := runtime.Caller(1)
	funcName := runtime.FuncForPC(pc).Name() // 获取函数调用者的名字
	log.Println("Function: " + funcName + "  message: \n" + fmt.Sprint(smg))
}

// TimeCost 功能：耗时统计 使用方式：在行数首行执行 defer TimeCost()()
func TimeCost() func() {
	start := time.Now()
	pc, _, _, _ := runtime.Caller(1)
	funcName := runtime.FuncForPC(pc).Name() // 获取函数调用者的名字
	return func() {
		tc := time.Since(start)
		log.Printf("Function: "+funcName+"(). Run finished. Time cost = %v\n", tc)
		// fmt.Printf("time cost = %v\n", tc)
	}
}
