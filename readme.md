# Selfhttps
Start a https proxy server with self-signed certificate.

## USAGE
Usage:

    selfhttps [-p PORT] -d DOMAIN1::PROXY_PASS1 [-d DOMAIN2::PROXY_PASS2] ...

Example:

    $ selfhttps - selfhttps -p 4430 -d local1.com::http://upstream1:4500 -d local2.com::http://upstream2:4501

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
                   
## Add trusted certificate(optional)
Add trusted certificate to system: 

    sudo security add-trusted-cert -d -r trustRoot -k /Library/Keychains/System.keychain ~/.selfhttps/local1.com.crt 

Remove trusted certificate from system: 

    sudo security delete-certificate -t -c local1.com 

## Change Log
- [x] Support websocket over https proxy
- [] Suport Yaml
