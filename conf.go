package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"regexp"
	"strings"
)

type DomainProxy struct {
	Domain    string
	ProxyPass string
}

type Config struct {
	Port string
	Silent bool
	DomainProxys []DomainProxy
	// Domain      string
	// ProxyPass   string
	// CertPath    string
	// CertKeyPath string
}

var conf *Config

func initConf(conf *Config, domainProxys []string) {
	if len(domainProxys) == 0 {
		fmt.Println("Usage: selfhttps -d local.com=http://localhost:5000")
		os.Exit(0)
	}
	for _, domainProxy := range domainProxys {
		domainProxyArr := strings.Split(strings.TrimSpace(domainProxy), "=")
		if len(domainProxyArr) != 2 {
			log.Fatalf("Invalid domain=proxy_pass(%s)\n", domainProxy)
		} else {
			domain, proxyPass := domainProxyArr[0], domainProxyArr[1]
			domain = strings.ToLower(domain)
			if !regexp.MustCompile(`^\w[\w\-\.]*$`).MatchString(domain) {
				fmt.Println("Invalid domain: " + domain)
				os.Exit(1)
			}
			u, err := url.Parse(proxyPass)
			if err != nil || u.Scheme == "" || domain == "" {
				fmt.Println("Usage: selfhttps -d local.com=http://10.120.45.10:5000")
				os.Exit(1)
			}
			conf.DomainProxys = append(conf.DomainProxys, DomainProxy{
				Domain:    domain,
				ProxyPass: proxyPass,
			})
		}
	}

	initCert(conf)
}
func getCertPath(domain string) (certPath, certKeyPath string) {
	home := os.Getenv("HOME")
	certPath = fmt.Sprintf(home+"/.selfhttps/%s.crt", domain)
	certKeyPath = fmt.Sprintf(home+"/.selfhttps/%s.key", domain)
	return certPath, certKeyPath
}

func initCert(conf *Config) {
	for _, domainProxy := range conf.DomainProxys {
		domain := domainProxy.Domain
		proxyPass := domainProxy.ProxyPass
		confCertPath, confCertKeyPath := getCertPath(domain)
		_, err1 := os.Stat(confCertPath)
		_, err2 := os.Stat(confCertKeyPath)
		// 0. generate certificate
		if os.IsNotExist(err1) || os.IsNotExist(err2) {
			cmd := fmt.Sprintf(`openssl req -x509 -nodes -newkey rsa:2048 -days 365 -keyout %s -out %s -subj "/C=CN/ST=GD/L=SZ/O=SelfHttps, Inc./CN=%s" -addext "subjectAltName = DNS:%s"`, confCertKeyPath, confCertPath, domain, domain)
			fmt.Printf("Generate certificate: %s \n", cmd)
			out, err := RunCommand("sh", "-c", cmd)
			if err != nil {
				fmt.Printf("failed to execute cmd(%s), err: %v, stdout: %s\n\nPlease install it manually:\n\033[31m brew install openssl\033[0m\n", cmd, err, out)
				os.Exit(0)
			}

			// openssl req -x509 -nodes -newkey rsa:1024    -keyout nginx.key -out nginx.crt -days 3650 -subj "/C=CN/ST=Some-Province/O=Internet Widgets, Inc./CN=$DOMAIN" -addext "subjectAltName = DNS:$DOMAIN"
		} else if err1 != nil || err2 != nil {
			log.Fatalf("cert(%s) or key(%s) file exists, but stat failed, err1(%v), err2(%v)\n", confCertPath, confCertKeyPath, err1, err2)
			return
		}
		// 1. add trusted certificate to system
		cmd := fmt.Sprintf("sudo security add-trusted-cert -d -r trustRoot -k /Library/Keychains/System.keychain ~/.selfhttps/%s.crt", domain)
		if !isMacOsx(){
			cmd = fmt.Sprintf("sudo cp ~/.selfhttps/%s.crt /usr/local/share/ca-certificates/ && sudo update-ca-certificates", domain)
		}
		runCmdWithConfirm("Add trusted certificate to system", cmd, conf.Silent)

		// 2. remove trusted certificate from system
		cmd=fmt.Sprintf("sudo security delete-certificate -t -c %s ", domain)
		if !isMacOsx(){
			cmd = fmt.Sprintf("sudo rm /usr/local/share/ca-certificates/%s.crt && sudo update-ca-certificates", domain)
		}
		fmt.Printf("The way to remove trusted certificate from system: \n\033[32m  %s\033[0m\n\n", cmd)
		fmt.Printf("Have a try: \033[94m curl -v -k https://%s:%s \033[0m ", domain, conf.Port)
		fmt.Printf("(proxy_pass: \033[94m%s\033[0m)\n\n", proxyPass)
	}
}

func GetConfig() *Config {
	if conf == nil {
		conf = &Config{}
		home := os.Getenv("HOME")
		if !IsFileExists(home + "/.selfhttps") {
			err := os.Mkdir(home+"/.selfhttps", 0755)
			if err != nil {
				log.Fatalf("Failed to create directory: %v", err)
			}
		}
	}
	return conf
}
