package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"golang.org/x/net/proxy"
	"io"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

const (
	timeout = 5 * time.Second
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
		Chproxy := strings.ToLower(sc.Text())
		wg.Add(1)
		go func() {
			err := CheckProxySOCKS(Chproxy, &wg)
			if err != nil {
				return
			}
		}()
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
		//"https://api.techniknews.net/ipgeo/",
		"https://myipip.net/",
	}
	rand.Seed(time.Now().Unix())
	url := urls[rand.Intn(len(urls))]
	//response, err := httpClient.Get(url)
	request, err := http.NewRequest("GET", url, nil)
	//增加header选项
	request.Header.Set("Referer", url)
	request.Header.Set("User-Agent", "curl/7.77.0")
	request.Header.Set("Origin", "url")
	//处理返回结果

	response, err := httpClient.Do(request)

	if err != nil {
		return nil
	}

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return nil
	}
	//body, err := ioutil.ReadAll(response.Body)

	json.Unmarshal(body, &ip)
	sp := strings.Split(proxyy, ":")
	respIp := sp[0]
	port := sp[1]

	if ip.IP == respIp {
		fmt.Printf("%s:%s\n", respIp, port)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(response.Body)
	return nil
}
