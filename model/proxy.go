package model

const (
	ProxyConfig = "server {\n listen  80;\n server_name %s;\n location / {\n proxy_pass  http://localhost:%d;\n}\n}\n#<?>\n"
	Placeholder = "#<?>"
)
