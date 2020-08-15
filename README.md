# bdy.get

# 用法
```
  PHP: url/?d=shareid&p=password
  GO: ./bdy_darwin_amd64 -d ShareID -p Password
```

# 参数
```
GO:
  -B BDUSS
    	Cookie:BDUSS：在baidu.com下
  -S STOKEN
    	Cookie:STOKEN：在pan.baidu.com下
  -c	创建或修改BDUSS和STOKEN
  -d ShareID
    	ShareID：分享ID
  -h	显示帮助
  -p Password
    	Password：提取码
  -t	测试BDUSS和STOKEN
  -v	显示版本
当参数-B,-S使用时，优先使用参数的数据。下载文件需要设置User-Agent:LogStatistic
```

```
PHP:
  -B BDUSS
    	Cookie:BDUSS：在baidu.com下
  -S STOKEN
    	Cookie:STOKEN：在pan.baidu.com下
  -c	创建或修改BDUSS和STOKEN
  -d ShareID
    	ShareID：分享ID
  -h	显示帮助
  -p Password
    	Password：提取码
  -t	测试BDUSS和STOKEN
  -v	显示版本
参数优先级h>v>c>t。下载文件需要设置User-Agent:LogStatistic
```


# Thanks

[PanDownload](https://pandownload.com)

[baiduwp](https://github.com/TkzcM/baiduwp)
