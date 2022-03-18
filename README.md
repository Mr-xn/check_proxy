# Soxy - a very fast tool for checking open SOCKS proxies in Golang 
I was looking for some open socks proxies, and so I needed to test them - but really fast. So I wrote on in Go!

### Installation
If you have a properly configured GOPATH and $GOPATH/bin is in your PATH, then run this command for a one-liner install, thank you golang!
```
go get -u github.com/pry0cc/soxy
```

### Usage
`proxies.txt`
```
8.8.8.8:3128
8.8.8.8:8080
```

```
cat proxies.txt | soxy | tee alive.txt
# if u want save pre results add -a after tee
cat proxies.txt | soxy | tee -a alive.txt
# test 9 times for filter the best proxies possible and use sort & uniq save result
# 1 from pbpaste (tested on mac  

file=socks.txt; for ((i=1;i<9;i++)) do pbpaste|./check_proxy_mac|tee -a $file; done && wc $file && sort $file|uniq|tee $file && wc $file

# 2 from file

file=socks.txt; for ((i=1;i<9;i++)) do cat $file|./check_proxy_mac|tee -a $file; done && wc $file && sort $file|uniq|tee $file && wc $file
```

### Credit
I pulled the proxy checking code and some of the multi-threading out of https://github.com/trigun117/ProxyChecker, so credit to trigun117! 
from https://github.com/pry0cc/soxy 
