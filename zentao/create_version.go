package main

import (
	// "net/url"
	"errors"
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var zentao_username,
	zentao_password,
	zentao_host_port,
	zentao_build_create_number, //40
	zentao_build_create_version,
	zentao_build_creator,
	zentao_build_product, //32
	zentao_build_branch, //0
	zentao_build_date,
	zentao_build_description string

func commandParamsParser() (err error) {
	oldUsage := flag.Usage
	flag.Usage = func() {
		oldUsage()
		fmt.Fprintf(os.Stdout, "\x0d\x0aIn topdfm, -cproduct=1 -cnumber=1 and none -cbranch"+
			"\x0d\x0aIn topmes, -cproduct=3 -cnumber=2 and none -cbranch"+
			"\x0d\x0aIn best-S1, -cproduct=32 -cnumber=40 and -cbranch=0"+
			"\x0d\x0a\x0d\x0aRelease date: 2018-10-30 +0800 wangxi"+
			"\x0d\x0a")
	}

	flag.StringVar(&zentao_username, "username", "", "zentao username. example: -username=mark or -username mark")
	flag.StringVar(&zentao_password, "password", "", "zentao password")
	flag.StringVar(&zentao_host_port, "hostPort", "http://192.196.1.2:9999", "zentao ip:port")
	flag.StringVar(&zentao_build_create_number, "cnumber", "", "zentao iteration number, see http://ip:port/zentao/build-create-${youriterationNo}.html to find out the appropriate number")
	flag.StringVar(&zentao_build_create_version, "cversion", "", "the new version that will be created")
	flag.StringVar(&zentao_build_creator, "creator", "", "the creator, normally, creator is equal to username")
	flag.StringVar(&zentao_build_product, "cproduct", "", "the No of product")
	flag.StringVar(&zentao_build_branch, "cbranch", "", "the No of branch")
	flag.StringVar(&zentao_build_date, "cdate", time.Now().Format("2006-01-02"), "creation date, default value is today")
	flag.StringVar(&zentao_build_description, "cdescription", "The version has been created automatically through auto-release", "creation description")
	flag.Parse()
	if 0 == len(zentao_username) ||
		0 == len(zentao_password) ||
		0 == len(zentao_build_create_number) ||
		0 == len(zentao_build_create_version) ||
		0 == len(zentao_build_product) {
		flag.Usage()
		err = errors.New("username, password, cnumber, cversion and cproduct both must not be empty. See -h for detail")
	}
	if 0 == len(zentao_build_creator) && 0 != len(zentao_username) {
		runes := []rune(zentao_username)
		zentao_build_creator = strings.ToUpper(string(runes[0:1])) + string(runes[1:])
	}
	return
}

func RandStringRunes(n int, params ...string) string {
	rand.Seed(time.Now().UnixNano())
	var letterRunes []rune
	if len(params) > 0 {
		letterRunes = []rune(params[0])
	} else {
		letterRunes = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789abcdefghijklmnopqrstuvwxyz")
	}
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func main() {

	//catch error
	// defer func() { //catch or finally
	//     if err := recover(); err != nil { //catch
	//         fmt.Fprintf(os.Stderr, "Exception: %v\n", err)
	//         os.Exit(1)
	//     }
	// }()

	err := commandParamsParser()
	if nil != err {
		fmt.Fprintf(os.Stderr, "Exception: %v\n", err)
		os.Exit(1)
		return
	}
	// os.Setenv("HTTP_PROXY", "http://127.0.0.1:8888")
	// proxyUrl, err := url.Parse("http://127.0.0.1:8888")
	// if nil != err {
	// 	fmt.Fprintf(os.Stderr, "Exception: %v\n", err)
	// 	return
	// }
	client := &http.Client{ /*Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)}*/ }
	//Login
	req, err := http.NewRequest("POST", zentao_host_port+"/zentao/user-login-L3plbnRhby9idWlsZC1jcmVhdGUtNDAuaHRtbA==.html",
		strings.NewReader("account="+zentao_username+"&password="+zentao_password+"&referer=%2Fzentao%2Fbuild-create-"+zentao_build_create_number+".html"))
	if nil != err {
		//to do, handle error
		fmt.Fprintf(os.Stderr, "Exception: %v\n", err)
		os.Exit(1)
		return
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64; rv:62.0) Gecko/20100101 Firefox/62.0")
	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Add("Accept-Language", "en-US,en;q=0.5")
	req.Header.Add("Accept-Encoding", "gzip, deflate")
	req.Header.Add("Referer", zentao_host_port+"/zentao/user-login-L3plbnRhby9idWlsZC1jcmVhdGUtNDAuaHRtbA==.html")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len("account="+zentao_username+"&password="+zentao_password+"&referer=%2Fzentao%2Fbuild-create-"+zentao_build_create_number+".html")))
	//Note: zentaosid must be lowercase and number
	zentaosid := RandStringRunes(26, "0123456789abcdefghijklmnopqrstuvwxyz")
	req.Header.Add("Cookie", "lang=zh-cn; device=desktop; theme=default; windowWidth=1525; windowHeight=848; zentaosid="+zentaosid)
	req.Header.Add("Upgrade-Insecure-Requests", "1")

	resp, err := client.Do(req)
	if err != nil {
		//to do, handle error
		fmt.Fprintf(os.Stderr, "Exception: %v\n", err)
		os.Exit(1)
		return
	}
	resp.Body.Close()

	//Create versioin
	// versionName := string("test_auto_release_1")
	// creator := string("Mark")
	// creationDate := time.Now().Format("2006-01-02")
	creationUid := RandStringRunes(13, "0123456789abcdefghijklmnopqrstuvwxyz")
	uniqueSqe := RandStringRunes(13, "123456789")
	postBody := string("")
	if 0 != len(zentao_build_product) {
		postBody += "-----------------------------${unique_seq}\x0d\x0aContent-Disposition: form-data; name=\"product\"\x0d\x0a\x0d\x0a${product}\x0d\x0a"
	}
	if 0 != len(zentao_build_branch) {
		postBody += "-----------------------------${unique_seq}\x0d\x0aContent-Disposition: form-data; name=\"branch\"\x0d\x0a\x0d\x0a${branch}\x0d\x0a"
	}
	if 0 != len(zentao_build_create_version) {
		postBody += "-----------------------------${unique_seq}\x0d\x0aContent-Disposition: form-data; name=\"name\"\x0d\x0a\x0d\x0a${version_name}\x0d\x0a"
	}
	if 0 != len(zentao_build_creator) {
		postBody += "-----------------------------${unique_seq}\x0d\x0aContent-Disposition: form-data; name=\"builder\"\x0d\x0a\x0d\x0a${creator}\x0d\x0a"
	}
	if 0 != len(zentao_build_date) {
		postBody += "-----------------------------${unique_seq}\x0d\x0aContent-Disposition: form-data; name=\"date\"\x0d\x0a\x0d\x0a${creation_date}\x0d\x0a"
	}
	postBody += "-----------------------------${unique_seq}\x0d\x0aContent-Disposition: form-data; name=\"scmPath\"\x0d\x0a\x0d\x0a\x0d\x0a"
	postBody += "-----------------------------${unique_seq}\x0d\x0aContent-Disposition: form-data; name=\"filePath\"\x0d\x0a\x0d\x0a\x0d\x0a"
	postBody += "-----------------------------${unique_seq}\x0d\x0aContent-Disposition: form-data; name=\"files[]\"; filename=\"\"\x0d\x0aContent-Type: application/octet-stream\x0d\x0a\x0d\x0a\x0d\x0a"
	postBody += "-----------------------------${unique_seq}\x0d\x0aContent-Disposition: form-data; name=\"labels[]\"\x0d\x0a\x0d\x0a\x0d\x0a"
	if 0 != len(zentao_build_description) {
		postBody += "-----------------------------${unique_seq}\x0d\x0aContent-Disposition: form-data; name=\"desc\"\x0d\x0a\x0d\x0a${description}\x0d\x0a"
	}
	if 0 != len(creationUid) {
		postBody += "-----------------------------${unique_seq}\x0d\x0aContent-Disposition: form-data; name=\"uid\"\x0d\x0a\x0d\x0a${creation_uid}\x0d\x0a"
	}
	postBody += "-----------------------------${unique_seq}--\x0d\x0a"
	postBody = strings.Replace(postBody, "${product}", zentao_build_product, -1)
	postBody = strings.Replace(postBody, "${branch}", zentao_build_branch, -1)
	postBody = strings.Replace(postBody, "${version_name}", zentao_build_create_version, -1)
	postBody = strings.Replace(postBody, "${creator}", zentao_build_creator, -1)
	postBody = strings.Replace(postBody, "${creation_date}", zentao_build_date, -1)
	postBody = strings.Replace(postBody, "${creation_uid}", creationUid, -1)
	postBody = strings.Replace(postBody, "${unique_seq}", uniqueSqe, -1)
	postBody = strings.Replace(postBody, "${description}", zentao_build_description, -1)

	req, err = http.NewRequest("POST", strings.Replace(zentao_host_port+"/zentao/build-create-${zentao_build_create_number}.html", "${zentao_build_create_number}", zentao_build_create_number, -1), strings.NewReader(postBody))
	if nil != err {
		//to do, handle error
		fmt.Fprintf(os.Stderr, "Exception: %v\n", err)
		os.Exit(1)
		return
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64; rv:62.0) Gecko/20100101 Firefox/62.0")
	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Add("Accept-Language", "en-US,en;q=0.5")
	req.Header.Add("Accept-Encoding", "gzip, deflate")
	req.Header.Add("Referer", strings.Replace(zentao_host_port+"/zentao/build-create-${zentao_build_create_number}.html", "${zentao_build_create_number}", zentao_build_create_number, -1))
	req.Header.Add("Content-Type", "multipart/form-data; boundary=---------------------------"+uniqueSqe)
	req.Header.Add("Content-Length", strconv.Itoa(len(postBody)))
	req.Header.Add("Cookie", "lang=zh-cn; device=desktop; theme=default; windowWidth=1525; windowHeight=848; zentaosid="+zentaosid)
	req.Header.Add("Upgrade-Insecure-Requests", "1")

	resp, err = client.Do(req)
	if err != nil {
		//to do, handle error
		fmt.Fprintf(os.Stderr, "Exception: %v\n", err)
		os.Exit(1)
		return
	}
	resp.Body.Close()

	fmt.Println("version has been created")
	os.Exit(0)
}
