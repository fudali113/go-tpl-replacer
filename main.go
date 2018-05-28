package main

import (
	"bytes"
	"flag"
	"fmt"
	"html/template"
	"io"
	"log"
	"os"
	"strings"
	"bufio"
)

func main() {

	var tplFiles, args, argFiles, outSplitStr string
	flag.StringVar(&tplFiles, "tpl-files", "config.tpl", "需要替換的模板文件; 多個文件使用 `;` 分割")
	flag.StringVar(&outSplitStr, "out-split", "---", "多個輸出的分割")
	flag.StringVar(&args, "args", "", "參數集合； 參數使用=分割, k=v ; 多個參數使用 `;` 分割")
	flag.StringVar(&argFiles, "arg-files", "", "參數文件集合, 多個文件使用 `;` 分割； 目前只支持 properties 文件格式")
	flag.Parse()

	replace, err := template.ParseFiles(strings.Split(tplFiles, ";")...)
	if err != nil {
		log.Printf("加載 tpl 文件出錯, error: %s", err.Error())
		os.Exit(1)
	}
	argsMap := map[string]interface{}{}
	if argFiles != "" {
		loadArgsByFile(argsMap, argFiles)
	}
	loadArgs(argsMap, args)
	replaces := replace.Templates()
	for i := range replaces {
		err = replaces[i].Execute(os.Stdout, argsMap)
		if err != nil {
			log.Printf("替換配置參數出錯, name: %s ; error: %s", replaces[i].Name(), err.Error())
		} else {
			io.Copy(os.Stdout, bytes.NewBufferString(fmt.Sprintf(" \n\n%s\n\n", outSplitStr)))
		}
	}
	os.Exit(0)

}

// loadArgs 解析命令行參數
func loadArgs(context map[string]interface{}, argsStr string) {
	args := strings.Split(argsStr, ";")
	for _, arg := range args {
		loadKvString(context, arg)
	}
}

// loadArgsByFile 解析文件參數
func loadArgsByFile(context map[string]interface{}, argFiles string) {
	filePaths := strings.Split(argFiles, ";")
	for _, filePath := range filePaths {
		func() {
			f, err := os.Open(filePath)
			if err != nil {
				log.Printf("打開文件 %s 出錯, err: %s ", filePath, err.Error())
				return
			}
			defer f.Close()
			rd := bufio.NewReader(f)
			for  {
				line, err := rd.ReadString('\n')
				loadKvString(context, line)
				if err != nil {
					break
				}
			}
		}()
	}
}

// loadKvString 加載 kv String 到參數列表
func loadKvString(context map[string]interface{}, kvString string)  {
	kv := strings.SplitN(kvString, "=", 2)
	if len(kv) != 2 {
		log.Printf("參數 %s 不能被正確解析", kvString)
	} else {
		context[kv[0]] = kv[1]
	}
}
