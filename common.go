package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"time"
)

func errorExit(err error) {
	if err != nil {
		pc, _, _, _ := runtime.Caller(1)
		funcName := runtime.FuncForPC(pc).Name() // 获取函数调用者的名字
		log.Println("-----------------------------\nError, exit.\n Function: " + funcName + "Error:\n" + error.Error(err))
		os.Exit(1)
	}
}

func errorTolerate(err error) {
	if err != nil {
		pc, _, _, _ := runtime.Caller(1)
		funcName := runtime.FuncForPC(pc).Name() // 获取函数调用者的名字
		log.Println("Error, but tolerate it.\n Function: " + funcName + "Error:\n" + error.Error(err))
	}
}
func printLog(smg interface{}) {
	pc, _, _, _ := runtime.Caller(1)
	funcName := runtime.FuncForPC(pc).Name() // 获取函数调用者的名字
	log.Println("Function: " + funcName + "  message: \n" + fmt.Sprint(smg))
}

// 功能：耗时统计 使用方式：在行数首行执行 defer timeCost()()
func timeCost() func() {
	start := time.Now()
	pc, _, _, _ := runtime.Caller(1)
	funcName := runtime.FuncForPC(pc).Name() // 获取函数调用者的名字
	return func() {
		tc := time.Since(start)
		log.Printf("Function: "+funcName+"(). Run finished. Time cost = %v\n", tc)
		// fmt.Printf("time cost = %v\n", tc)
	}
}
