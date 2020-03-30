package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"html/template"
	"io"
	"log"
	"os"
	"strings"
)

func main() {

	var tplFiles, args, argFiles, outSplitStr, getArgValue string
	flag.StringVar(&tplFiles, "tpl-files", "config.tpl", "需要替換的模板文件; 多個文件使用 `;` 分割")
	flag.StringVar(&outSplitStr, "out-split", "---", "多個輸出的分割")
	flag.StringVar(&args, "args", "", "參數集合； 參數使用=分割, k=v ; 多個參數使用 `;` 分割")
	flag.StringVar(&argFiles, "arg-files", "", "參數文件集合, 多個文件使用 `;` 分割； 目前只支持 properties 文件格式")
	flag.StringVar(&getArgValue, "get-arg-value", "", "获取某个参数的值")
	flag.Parse()

	argsMap := map[string]interface{}{}
	if argFiles != "" {
		loadArgsByFile(argsMap, argFiles)
	}
	if args != "" {
		loadArgs(argsMap, args)
	}

	if getArgValue != "" {
		fmt.Printf("%s", argsMap[getArgValue])
		return
	}

	replace, err := template.ParseFiles(strings.Split(tplFiles, ";")...)
	if err != nil {
		log.Printf("加載 tpl 文件出錯, error: %s", err.Error())
		os.Exit(1)
	}
	replaces := replace.Templates()
	for i := range replaces {
		err = replaces[i].Execute(os.Stdout, argsMap)
		if err != nil {
			log.Printf("替換配置參數出錯, name: %s ; error: %s", replaces[i].Name(), err.Error())
			os.Exit(1)
		} else {
			io.Copy(os.Stdout, bytes.NewBufferString(fmt.Sprintf(" \n%s\n", outSplitStr)))
		}
	}
	os.Exit(0)

}

// loadArgs 解析命令行參數
func loadArgs(context map[string]interface{}, argsStr string) {
	args := strings.Split(argsStr, ";")
	for _, arg := range args {
		if arg == "" {
			continue
		}
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
				os.Exit(1)
			}
			defer f.Close()
			rd := bufio.NewReader(f)
			for  {
				line, err := rd.ReadString('\n')
				line = strings.TrimSpace(line)
				// 支持空格换行和 以 # 注释
				if line == "" || strings.HasPrefix(line, "#") {
					if io.EOF == err {
						break
					}
					continue
				}
				if err != nil && line == "" {
					log.Printf(err.Error())
					break
				}
				loadKvString(context, line)
			}
		}()
	}
}

// loadKvString 加載 kv String 到參數列表
func loadKvString(context map[string]interface{}, kvString string)  {
	kv := strings.SplitN(kvString, "=", 2)
	if len(kv) != 2 {
		log.Printf("參數 %s 不能被正確解析", kvString)
		os.Exit(1)
	} else {
		context[kv[0]] = kv[1]
	}
}
