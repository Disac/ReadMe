# TLS（Transport Layer Security：传输层安全协议）
协议分为两个部分：TLS记录协议和TLS握手协议。TLS用于通信之间数据的加密和认证。

## CA数字证书（包含CA的公钥）

1. 首先生成ca的key
```
$ openssl genrsa -out ca.key 2048
```

2. 根据自己CA的私钥自签发CA证书（证书中包含CA公钥）
```
$ openssl req -x509 -new -nodes -key ca.key -subj "/CN=dc.com" -days 3650 -out ca.crt
```

## Golang如何实现cli-server之间通过TLS加密（http层）

#### 1. 对服务端证书进行校验

1. 生成server端的私钥
```
$ openssl genrsa -out server.key 2048
```

2. 生成server证书的签名请求
```
$ openssl req -new -key server.key -subj "/CN=localhost" -out server.csr
```

3. 根据自己的CA使用CA的私钥对server的签名请求处理，得到server的数字证书
```
$ openssl x509 -req -in server.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out server.crt -days 3650
```

在我们的server中使用刚才生成的服务端证书：
```
http.ListenAndServeTLS(":8080", "server.crt", "server.key", nil)
```

在我们的cli中访问我们的server：
```
certPool := x509.NewCertPool()

caCrt, err := ioutil.ReadFile("ca.crt")

certPool.AppendCertsFromPEM(caCrt)

client := http.Client{Transport: &http.Transport{
	TLSClientConfig: &tls.Config{
		RootCAs: certPool,
	},
}}
client.Get("https://localhost:8080")
```

### 2. 对客户端证书进行校验
需要先为客户端生成证书
1. 生成client端的私钥
```
$ openssl genrsa -out client.key 2048
```

2. 生成server证书的签名请求
```
$ openssl req -new -key client.key -subj "/CN=localhost" -out client.csr
```

3. 根据自己的CA使用CA的私钥对client的签名请求处理，得到client的数字证书
```
$ openssl x509 -req -in client.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out client.crt -days 3650
```

server端需要校验client的数字证书，加载用于校验的ca.crt：
```
certPool := x509.NewCertPool()
caCrt, err := ioutil.ReadFile("ca.crt")

certPool.AppendCertsFromPEM(caCrt)

server := &http.Server{
	Addr:    ":8080",
	Handler: &TestHandler{},
	TLSConfig: &tls.Config{
		ClientCAs:  certPool,
		ClientAuth: tls.RequireAndVerifyClientCert,
	},
}
server.ListenAndServeTLS("server.crt", "server.key")
```

客户端也需要加载自身的数字证书，用于server端连接时做证书校验：
```
certPool := x509.NewCertPool()
caCrt, err := ioutil.ReadFile("ca.crt")

certPool.AppendCertsFromPEM(caCrt)

clientCrt, err := tls.LoadX509KeyPair("client.crt", "client.key")
client := http.Client{
	Transport: &http.Transport{
		TLSClientConfig: &tls.Config{
			RootCAs:      certPool,
			Certificates: []tls.Certificate{clientCrt},
		},
	},
}

client.Get("https://localhost:8080")
```

备注：可以查看自己的证书的cmd哦
```
$ openssl x509 -text -in xxx.crt -noout
```

## Golang实现ss通信通过TLS加密