package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/samber/lo"
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
)

func main() {
	conf := GetConfig()
	flags := []cli.Flag{
		&cli.StringSliceFlag{
			Name:     "d",
			Usage:    "domain and proxy_pass `DOMAIN::proxy_pass`",
			Required: true,
			// Destination: &conf.Domain,
		},
		&cli.StringFlag{
			Name:        "p",
			Value:       "443",
			Usage:       "server `PORT`",
			Destination: &conf.Port,
		},
		// &cli.StringFlag{
		// 	Name: "s",
		// 	// Aliases:     []string{"s"},
		// 	Value:       "http://localhost:4500",
		// 	Usage:       "Nginx's `proxy_pass`",
		// 	Destination: &conf.ProxyPass,
		// },
	}
	app := &cli.App{
		Name:        "selfhttps",
		Description: "start a https proxy server with self-signed certificate",
		UsageText:   "selfhttps [-p PORT] -d DOMAIN1::PROXY_PASS1 [-d DOMAIN2::PROXY_PASS2] ...",
		Usage: `selfhttps -p 4430 -d local1.com::http://upstream1:4500 -d local2.com::http://upstream2:4501

echo "127.0.0.1 local1.com local2.com upstream1 upstream2" | sudo tee -a /etc/hosts
curl -v -k https://local1.com:4430/api/v1/xxx
curl -v -k https://local2.com:4430/api/v1/xxx

        +----------------+
        | curl/Chrome/...|
        +------+---------+
               |
               v 
		   +-------+------+
		   | https proxy  | default port: 443
		   | (port:4430)  |  
		   ++-----+-------+  
          |         | (like nginx's proxy_pass)
          v         v
+-------+---+        +-----------+  
| upstream1 |        | upstream2 |  
|(port:4500)|        |(port:4501)|  
+-----------+        +-----------+  
		`,
		Before: altsrc.InitInputSourceWithContext(flags, altsrc.NewYamlSourceFromFlagFunc("load")),
		Flags:  flags,
		Action: func(cCtx *cli.Context) error {
			domainProxys := cCtx.StringSlice("d")
			initConf(conf, domainProxys)
			quit := make(chan os.Signal, 1)
			proxyServer := createProxyServer(func() {
				log.Println("cleanup")
			})
			domains := lo.Map(conf.DomainProxys, func(dp DomainProxy, i int) string {
				return dp.Domain
			})

			fmt.Printf("Config hosts:\n\techo '127.0.0.1 %s' | sudo tee -a /etc/hosts\n", strings.Join(domains, " "))
			fmt.Printf("Press Ctrl+C to shutdown\n")
			signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
			sig := <-quit
			err := proxyServer.Close()
			log.Printf("gracefully shutdown Server ...(sig=%v, err=%v)\n", sig, err)
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}

}
