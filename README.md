## go-tpl-replacer

使用 go template 模板對文件進行配置並執行輸出數據到 標準輸出流

### usage

such as:
```
./go-tpl-replacer -tpl-files config1.tpl -arg-files args.properties -out-split --- -args name=test name

out:
hello:
    name: test
    value: i am value

---
```

### help

```
./go-tpl-replacer -help

Usage of ./go-tpl-replacer:
  -arg-files ;
        參數文件集合, 多個文件使用 ; 分割； 目前只支持 properties 文件格式
  -args ;
        參數集合； 參數使用=分割, k=v ; 多個參數使用 ; 分割
  -out-split string
        多個輸出的分割 (default "---")
  -tpl-files ;
        需要替換的模板文件; 多個文件使用 ; 分割 (default "config.tpl")


```

