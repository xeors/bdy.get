package main

import (
	"io/ioutil"
	"net/http"
	"strings"
	"regexp"
	"os"
	"fmt"
	"flag"
	json "github.com/tidwall/gjson"
)

func usage() {
    print("\n" + appname + "/" + version + "\n用法：" + os.Args[0] + " -d ShareID -p Password\n参数：\n")
	flag.PrintDefaults()
	print("当参数-B,-S使用时，优先使用参数的数据。下载文件需要设置User-Agent:LogStatistic\n\n")
}

var appname string = "bdy.get"
var version string = "1.0.0g"

func main()  {
	B := ""
	S := ""
	path := "bdy.inf"
	var h, t, c, v bool
	var d, p string
	flag.StringVar(&d, "d", "", "`ShareID`：分享ID")
	flag.StringVar(&p, "p", "", "`Password`：提取码")
	flag.StringVar(&B, "B", "", "Cookie:`BDUSS`：在baidu.com下")
	flag.StringVar(&S, "S", "", "Cookie:`STOKEN`：在pan.baidu.com下")
	flag.BoolVar(&t, "t", false, "测试BDUSS和STOKEN")
	flag.BoolVar(&c, "c", false, "创建或修改BDUSS和STOKEN")
	flag.BoolVar(&v, "v", false, "显示版本")
	flag.BoolVar(&h, "h", false, "显示帮助")
	flag.Parse()
	if h {
		flag.Usage = usage
		flag.Usage()
		os.Exit(0)
	}
	if v {
		print(appname + "/" + version + "\n")
		os.Exit(0)
	}
	if c {
		if B == "" || S == "" {
			fmt.Print("BDUSS:")
			fmt.Scan(&B)
			fmt.Print("STOKEN:")
			fmt.Scan(&S)
		}
		w, err := os.Create(path)
		defer w.Close()
		if err !=nil {
			print("文件创建失败!!\n")
			os.Exit(0)
		} else {
			_, err=w.Write([]byte(B + ";;" + S))
			if err != nil {
				print("写入失败!!\n")
			}
		}
		os.Exit(0)
	}
	if B == "" || S == "" {
		if isex(path) {
			r, err := os.OpenFile(path, os.O_RDONLY, 0600)
			defer r.Close()
			if err !=nil {
				print("文件打开失败!!\n")
				os.Exit(0)
			}else{
				content, err := ioutil.ReadAll(r)
				if err == nil {
					R := string(content)
					if len(R) > 2 {
						BS := strings.Split(R, ";;")
						if len(BS) == 2 {
							B = BS[0]
							S = BS[1]
						}else{
							print("数据错误!!\n")
							os.Exit(0)
						}
					}else{
						print("数据错误!!\n")
						os.Exit(0)
					}
				}else{
					print("读取错误!!\n")
					os.Exit(0)
				}
			}
		}else{
			fmt.Print("首次使用需要配置信息，请输入以下信息：\nBDUSS:")
			fmt.Scan(&B)
			fmt.Print("STOKEN:")
			fmt.Scan(&S)
			w, err := os.Create(path)
			defer w.Close()
			if err !=nil {
				print("文件创建失败!!\n")
				os.Exit(0)
			} else {
				_, err=w.Write([]byte(B + ";;" + S))
				if err != nil {
					print("写入失败!!\n")
				}
			}
		}		
	}
	var R1, R2, R3, R4, R5 Rr
	STOKEN := S
	BDUSS := B
	if t {
		var test Rr
		test.u = "https://pan.baidu.com/disk/home"
		test.h = [][]string{
			{"Cookie", "STOKEN=" + STOKEN + ";BDUSS=" + BDUSS},
		}
		testr := Http(test)
		if len(testr.b) < 23 {
			print("登陆信息无效!!\n")
		}else{
			print("登陆信息有效!!\n")
		}
		print(len(testr.b))
		os.Exit(0)
	}
	if d == "" || p == "" {
		if len(os.Args) == 3 {
			d = os.Args[1]
			p = os.Args[2]
		}else if len(os.Args) == 2 && strings.Index(os.Args[1], "提取码") != -1 {
			d = match(os.Args[1], `baidu.com\/s\/(.*?) `, "ERROR!!0x00")[0][1]
			p = match(os.Args[1], `码: (.*?)$`, "ERROR!!0x01")[0][1]
		}else{
			fmt.Print("请输入分享ID：")
			fmt.Scan(&d)
			fmt.Print("请输入提取码：")
			fmt.Scan(&p)
		}
	}
	if d != "" && p != "" {
		if d[0:1] == "1" {
			d = d[1:len(d)]
		}
		R1.u = "https://pan.baidu.com/share/verify?surl=" + d
		R1.d = "pwd=" + p
		R1.h = [][]string{
			{"User-Agent", "netdisk"},
			{"Referer", "https://pan.baidu.com/disk/home"},
		}
		r1 := Http(R1)
		errno := json.Get(r1.b, "errno").String()
		if errno == "0" {
			randsk := json.Get(r1.b, "randsk").String()
			R2.u = "https://pan.baidu.com/s/1" + d
			h := [][]string{
				{"Cookie", "STOKEN=" + STOKEN + ";BDUSS=" + BDUSS + ";BDCLND=" + randsk},
				{"User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.514.1919.810 Safari/537.36"},
			}
			R2.h = h
			r2 := Http(R2)
			m2 := match(r2.b, `yunData.setData\(({.+)\);`, "ERROR!!0x02")[0][1]
			sign := json.Get(m2, "sign").String()
			timestamp := json.Get(m2, "timestamp").String()
			shareid := json.Get(m2, "shareid").String()
			uk := json.Get(m2, "uk").String()
			R3.u = "https://pan.baidu.com/share/list?app_id=250528&channel=chunlei&clienttype=0&desc=0&num=100&order=name&page=1&root=1&shareid=" + shareid + "&showempty=0&uk=" + uk + "&web=1"
			R3.h = h
			r3 := Http(R3)
			errno = json.Get(r3.b, "errno").String()
			if errno == "0" {
				fsid := json.Get(r3.b, "list.0.fs_id").String()
				R4.u = "https://pan.baidu.com/api/sharedownload?app_id=250528&channel=chunlei&clienttype=12&sign=" + sign + "&timestamp=" + timestamp + "&web=1"
				R4.d = `encrypt=0&extra={"sekey":"` + randsk + `"}&fid_list=[` + fsid + "]&primaryid=" + shareid + "&uk=" + uk + "&product=share&type=nolimit"
				R4.h = h
				r4 := Http(R4)
				errno = json.Get(r4.b, "errno").String()
				if errno == "0" {
					R5.u = json.Get(r4.b, "list.0.dlink").String()
					R5.h = [][]string{
						{"User-Agent", "LogStatistic"},
						{"Cookie", "BDUSS=" + BDUSS},
					}
					r5 := Http(R5)
					m5 := match(r5.rh, `Location:(.*?)\r\n`, "ERROR!!0x03")[0][1]
					print(m5 + "\n")
					if len(os.Args) == 1 {
						fmt.Print("按下回车键退出...")
						fmt.Scanln(&d)
					}
					os.Exit(0)
				}
			}
		}else if(errno == "-12"){
			print("提取码错误或文件失效!!\n")
		}else{
			print("ERROR!!0x04\n")
		}
	}
}

type Rr struct {
	Url,u string
	Header,h [][]string
	Status,s string
	Data,d string
	Body,b string
	Method,m string
	ReturnHeader,rh string
}

func Http(Req Rr) Rr{
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		 },
	}
	D := strings.NewReader(Req.d)
	if Req.d != "" && Req.m == "" {
		Req.m = "POST"
	}
	request, err := http.NewRequest(Req.m, Req.u, D)
	if(len(Req.h) != 0){
        request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		for i := 0; i < len(Req.h); i++ {
			request.Header.Set(Req.h[i][0], Req.h[i][1])
        }
	}
	if err != nil {
		print("ERROR!!0x22\n")
		os.Exit(0)
	}
    response, err := client.Do(request)
    if err != nil {
        print("ERROR!!0x23\n")
        os.Exit(0)
    }
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		print("ERROR!!0x33\n")
		os.Exit(0)
	}
	h := response.Header
	hs := ""
	for k := range h {
        hs = hs + k + ":" + h[k][0] + "\r\n"
	}
	var ret Rr
	ret.rh = hs + response.Status + "\r\n"
	ret.b = string(body)
	ret.s = response.Status
	return ret
}

func match(nr, reg, err string) [][]string{
	p := regexp.MustCompile(reg)
	result := p.FindAllStringSubmatch(nr,-1)
	if len(result) != 0 {
		return result
	}else{
        if err != "" {
            print(err + "\n")
            os.Exit(0)
        }
        return nil
	}
}

func isex(path string) (bool) {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}