# Proxy Go + Cloudflare Durable Containers + Workers

Projeto completo contendo:

- Proxy Go customizado
- Container Cloudflare
- Worker roteando tráfego
- Load balance
- Containers por ID

## 🛠 Como usar

### 1) Clone

git clone https://github.com/seu-usuario/meu-proxy-container.git

### 2) Entre na pasta

cd meu-proxy-container

### 3) Publique a imagem do container

cd container docker build -t proxycontainer . wrangler publish

### 4) Deploy do Worker

cd ../worker npm install wrangler deploy

### 5) Teste

https://SEU-WORKER.workers.dev/container/teste
https://SEU-WORKER.workers.dev/lb
https://SEU-WORKER.workers.dev/singleton

