# Selfhttps
Start a https proxy server with self-signed certificate.

- [x] Websocket over https proxy
- [x] Auto generated certificate
- [x] Support MacOSX, linux, windows(partial)

## install

    bash -c "$(curl -fsSL https://raw.githubusercontent.com/ahuigo/selfhttps/main/install.sh)"

## USAGE
Usage:
    
    $ selfhttps -h
    selfhttps [-p PORT] [--silent] -d domain1=proxy_pass1 [-d domain2=proxy_pass2] ...

Example:

    $ selfhttps - selfhttps -d local1.com=http://upstream1:4500 -d local2.com=http://upstream2:4501

    echo "127.0.0.1 local1.com local2.com upstream1 upstream2" | sudo tee -a /etc/hosts
    curl -v -k https://local1.com/api/v1/xxx
    curl -v -k https://local2.com/api/v1/xxx

                    +---------------------------+
                    |curl -k https://local1.com |
                    |curl -k https://local2.com |
                    +------+--------------------+
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
                   
## Add trusted certificate to OS(to ignore cert warnning)
> If you don't wanna see certificate warnning, you could put certificate into your OS system.

Add trusted certificate to system: 

    # mac
    sudo security add-trusted-cert -d -r trustRoot -k /Library/Keychains/System.keychain ~/.selfhttps/local1.com.crt 

    # linux(ubuntu/debian)
    sudo cp ~/.selfhttps/local1.com.crt /usr/local/share/ca-certificates/ && sudo update-ca-certificates

    # windows
    certutil -addstore -f "ROOT" /path/to/.selfhttps/local1.com.crt

Remove trusted certificate from system: 

    # mac
    sudo security delete-certificate -t -c local1.com 

    # linux
    sudo rm /usr/local/share/ca-certificates/local1.com.crt && sudo update-ca-certificates

    # windows
    certutil -delstore "ROOT" /path/to/.selfhttps/local1.com.crt
