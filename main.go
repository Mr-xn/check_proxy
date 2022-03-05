package main

import (
	"bufio"
	"encoding/json"
	"fmt"
        "math/rand"
	"golang.org/x/net/proxy"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

const (
	timeout = time.Duration(5 * time.Second)
)

// IP is a struct for storing the JSON output from apify
type IP struct {
	IP string
}

var wg sync.WaitGroup

func main() {
	runtime.GOMAXPROCS(50)

	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		proxy := strings.ToLower(sc.Text())
		wg.Add(1)
		go CheckProxySOCKS(proxy, &wg)
	}

	wg.Wait()
}

//CheckProxySOCKS Check proxies on valid
func CheckProxySOCKS(proxyy string, wg *sync.WaitGroup) (err error) {

	defer wg.Done()

	d := net.Dialer{
		Timeout:   timeout,
		KeepAlive: timeout,
	}

	//Sending request through proxy
	dialer, _ := proxy.SOCKS5("tcp", proxyy, nil, &d)
	var ip IP
	httpClient := &http.Client{
		Timeout: timeout,
		Transport: &http.Transport{
			DisableKeepAlives: true,
			Dial:              dialer.Dial,
		},
	}

	urls := []string{
		"https://api-ipv4.ip.sb/jsonip",
		"https://api.ipify.org?format=json",
		"https://ip-fast.com/api/ip/?format=json",
		"https://api.myip.com/",
		"https://ipinfo.io/widget",
		"https://ipapi.co/json",
		"https://api.techniknews.net/ipgeo/",
	}
	rand.Seed(time.Now().Unix())
	url := urls[rand.Intn(len(urls))]
	//response, err := httpClient.Get("https://api-ipv4.ip.sb/jsonip")
	request, err := http.NewRequest("GET", url, nil)
	//增加header选项
	request.Header.Add("Referer", url)
	request.Header.Add("User-Agent", "curl/7.77.0")
	request.Header.Add("Origin", "url")

	if err != nil {
		return
	}

	//处理返回结果
	response, _ := httpClient.Do(request)
	body, err := ioutil.ReadAll(response.Body)
	//body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return
	}

	json.Unmarshal([]byte(body), &ip)
	sp := strings.Split(proxyy, ":")
	respIp := sp[0]
	port := sp[1]

	if ip.IP == respIp {
		fmt.Printf("%s:%s\n", respIp, port)
	}
	defer response.Body.Close()
	return nil
}
