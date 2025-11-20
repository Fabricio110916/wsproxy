package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	TargetAddr = "137.131.138.10"
)

const (
	ServerPort       = "8080"           // Cloudflare exige HTTP na porta 8080
	TargetPortSSH    = "22"
	TargetPortV2Ray  = "8282"
	BufferSize       = 524288           // 512 KB
	KeepAliveTimeout = 24 * time.Hour   // Mantém conexão viva por até 24h
)

type Target struct {
	Addr  string
	Port  string
	V2Ray bool
}

func createTarget(endpoint string) *Target {
	if endpoint == "/ws/" {
		return &Target{Addr: TargetAddr, Port: TargetPortV2Ray, V2Ray: true}
	}
	return &Target{Addr: TargetAddr, Port: TargetPortSSH, V2Ray: false}
}

func copyStream(src, dst net.Conn, wg *sync.WaitGroup, direction string) {
	defer wg.Done()
	buffer := make([]byte, BufferSize)
	for {
		n, err := src.Read(buffer)
		if err != nil {
			if err != io.EOF {
				fmt.Printf("[ERROR] Copiando dados (%s): %v\n", direction, err)
			}
			break
		}
		if n > 0 {
			_, err := dst.Write(buffer[:n])
			if err != nil {
				fmt.Printf("[ERROR] Escrevendo dados (%s): %v\n", direction, err)
				break
			}
		}
	}
}

func keepAlive(conns ...net.Conn) {
	ticker := time.NewTicker(KeepAliveTimeout)
	defer ticker.Stop()

	for range ticker.C {
		for _, conn := range conns {
			conn.SetDeadline(time.Now().Add(KeepAliveTimeout))
		}
	}
}

// --- ADAPTADO PARA CLOUDLFRARE ----
// Transformar HTTP → TCP com Hijack
// ----------------------------------
func httpHandler(w http.ResponseWriter, r *http.Request) {
	hj, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "Hijack não suportado", http.StatusInternalServerError)
		return
	}

	client, rw, err := hj.Hijack()
	if err != nil {
		return
	}

	go handleClient(client, rw)
}

// --- SUA LÓGICA ORIGINAL (mantida totalmente) ----
func handleClient(client net.Conn, rw *bufio.ReadWriter) {
	defer client.Close()

	clientAddr := client.RemoteAddr().String()
	fmt.Printf("[INFO] Cliente conectado: %s\n", clientAddr)

	// Lê a primeira linha HTTP
	payload, err := rw.ReadString('\n')
	if err != nil {
		fmt.Printf("[ERROR] Falha lendo cabecalho: %v\n", err)
		return
	}

	parts := strings.Split(payload, " ")
	if len(parts) < 2 {
		return
	}

	endpoint := parts[1]
	target := createTarget(endpoint)

	targetConn, err := net.Dial("tcp", target.Addr+":"+target.Port)
	if err != nil {
		fmt.Printf("[ERROR] Falha ao conectar no destino: %v\n", err)
		return
	}
	defer targetConn.Close()

	if target.V2Ray {
		targetConn.Write([]byte(payload))
	} else {
		client.Write([]byte("HTTP/1.1 101 Switching Protocols\r\nUpgrade: websocket\r\nConnection: Upgrade\r\n\r\n"))
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go copyStream(client, targetConn, &wg, "Cliente → Alvo")
	go copyStream(targetConn, client, &wg, "Alvo → Cliente")

	go keepAlive(client, targetConn)

	wg.Wait()
	fmt.Printf("[INFO] Conexão encerrada: %s\n", clientAddr)
}

func main() {
	fmt.Printf("[INFO] Proxy rodando dentro de Container Cloudflare na porta %s\n", ServerPort)

	http.HandleFunc("/", httpHandler)

	err := http.ListenAndServe("0.0.0.0:"+ServerPort, nil)
	if err != nil {
		fmt.Printf("[FATAL] Erro ao iniciar servidor HTTP: %v\n", err)
		os.Exit(1)
	}
}