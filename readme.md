# Selfhttps
Start a https proxy server with self-signed certificate.

- [x] Websocket over https proxy
- [x] Auto generated certificate
- [ ] Suport Yaml

Required:
- go >= 1.22
- openssl >= 1.1.1 or LibreSSL >= 3.1.0
    - Mac OSX: brew install openssl
    - Debian/Ubuntu: sudo apt install openssl

## USAGE
Usage:

    selfhttps [-p PORT] -d domain1=proxy_pass1 [-d domain2=proxy_pass2] ...

Example:

    $ selfhttps - selfhttps -d local1.com=http://upstream1:4500 -d local2.com=http://upstream2:4501

    echo "127.0.0.1 local1.com local2.com upstream1 upstream2" | sudo tee -a /etc/hosts
    curl -v -k https://local1.com/api/v1/xxx
    curl -v -k https://local2.com/api/v1/xxx

                       +----------------+
                       | curl/Chrome/...|
                       +------+---------+
                              |
                              v 
                      +-------+------+
                      | https proxy  | default port: 443
                      | (port:443)  |  
                      ++-----+-------+  
                         |         | (same as nginx's proxy_pass)
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

