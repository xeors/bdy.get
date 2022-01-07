# bdy.get

# 用法
```
  ./bdy_darwin_amd64 -d ShareID -p Password
```

# 参数
```
bdy.get/1.1.3g
用法：bdy -d ShareID -p Password
参数：
  -B BDUSS
    	Cookie:BDUSS：在baidu.com下
  -C Command
    	Command：自定义输出，见下方说明
  -P Path
    	Path：配置文件，支持本地文件和远程文件(http)
  -S STOKEN
    	Cookie:STOKEN：在pan.baidu.com下
  -c	创建或修改配置文件
  -d ShareID
    	ShareID：分享ID
  -dir Directory
    	Directory: 当分享内容为文件夹时指定路径(0.0.2)
  -h	显示帮助
  -p Password
    	Password：提取码
  -s Server
    	Server: 启动为服务器模式，此时参数p作为指定端口，默认23333
  -t	测试BDUSS和STOKEN
  -v	显示版本
  -x Xsh
    	Xsh: Xshare
当参数-B,-S使用时，优先使用参数的数据。下载文件需要设置User-Agent:LogStatistic

参数C的常量和变量：
  ;;;  获取到的直链
  ;b;  BDUSS
  ;s;  STOKEN
  ;h;  Logtatistic(常量)
  ;ua; User-Agent:LogStatistic(常量)
例如：-C 'wget --header=";ua;" ";;;"'
```


# Thanks

[PanDownload](https://pandownload.com)

[baiduwp](https://github.com/TkzcM/baiduwp)
