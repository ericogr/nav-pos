package cmd

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strings"
	"syscall"

	v1 "github.com/ericogr/nav-pos/api/v1"
	app "github.com/ericogr/nav-pos/internal/app"
	"github.com/spf13/cobra"
)

//go:embed static/*
var staticFiles embed.FS

var (
	openBrowser       bool
	host              string
	port              int
	telemetryProvider string
	telemetryParams   string
	aircraftProvider  string
	aircraftParams    string
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the application",
	Run: func(cmd *cobra.Command, args []string) {
		start(openBrowser, host, port, telemetryProvider, telemetryParams, aircraftProvider, aircraftParams)
	},
}

func init() {
	rootCmd.AddCommand(startCmd)

	startCmd.Flags().BoolVarP(&openBrowser, "openBrowser", "o", true, "Open browser on start")
	startCmd.Flags().IntVarP(&port, "port", "p", 8080, "Port to bind the server to")
	startCmd.Flags().StringVar(&host, "host", "localhost", "Host to bind the server to")
	startCmd.Flags().StringVarP(&telemetryProvider, "tprovider", "t", "mavlinkserial", "Telemetry provider to use (ex fake, mavlinkserial)")
	startCmd.Flags().StringVar(&telemetryParams, "tparams", "", "Telemetry provider parameters")
	startCmd.Flags().StringVarP(&aircraftProvider, "aprovider", "a", "opensky", "Aircraft provider to use (ex fake, opensky)")
	startCmd.Flags().StringVar(&aircraftParams, "aparams", "", "Aircraft provider parameters")
}

func start(openBrowser bool, host string, port int, telemetryProvider string, telemetryParams string, aircraftProvider string, aircraftParams string) {
	ctx, cancel := context.WithCancel(context.Background())
	serverURL := fmt.Sprintf("http://%s:%d", host, port)

	telemetryParamsMap, err := paramsStringToMap(telemetryParams)
	if err != nil {
		log.Fatalf("Erro ao converter parâmetros de telemetria: %v", err)
	}
	telemetry, err := app.CreateTelemetry(telemetryProvider, telemetryParamsMap)
	if err != nil {
		log.Fatalf("Erro ao criar o provedor de telemetria: %v", err)
	}
	telemetry.Initialize(ctx)

	aircraftParamsMap, err := paramsStringToMap(aircraftParams)
	if err != nil {
		log.Fatalf("Erro ao converter parâmetros de aeronave: %v", err)
	}
	aircraft, err := app.CreateAircraft(aircraftProvider, aircraftParamsMap)
	if err != nil {
		log.Fatalf("Erro ao criar o provedor de aeronaves: %v", err)
	}

	handleRequests := v1.HandleRequests{
		Telemetry: telemetry,
		Aircraft:  aircraft,
	}
	http.HandleFunc("/api/v1/aircraft", handleRequests.HandleAircraft)
	http.HandleFunc("/api/v1/telemetry", handleRequests.HandleTelemetry)
	fSys, _ := fs.Sub(staticFiles, "static")
	http.Handle("/", http.FileServer(http.FS(fSys)))

	// Iniciar o servidor
	go func() {
		log.Printf("Server started at %s", serverURL)
		err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
		if err != nil {
			log.Printf("Server error: %v", err)
		}
	}()

	err = openUrlOnBrowser(serverURL)
	if err != nil {
		log.Printf("Error opening browser: %v", err)
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	sig := <-stop
	log.Printf("\nReceived signal: %s\n", sig)
	cancel()
	log.Println("Terminating server...")
}

// Função para abrir o navegador
func openUrlOnBrowser(url string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	case "darwin":
		cmd = exec.Command("open", url)
	default: // "linux", "freebsd", etc.
		cmd = exec.Command("xdg-open", url)
	}

	return cmd.Start()
}

func paramsStringToMap(params string) (map[string]string, error) {
	// Implementar a lógica para converter a string de parâmetros em um mapa
	// Exemplo: "param1=value1,param2=value2" -> map[string]string{"param1": "value1", "param2": "value2"}
	paramsMap := make(map[string]string)
	if len(params) > 0 {
		pairs := strings.Split(params, ",")
		for _, pair := range pairs {
			kv := strings.Split(pair, "=")
			if len(kv) != 2 {
				return nil, fmt.Errorf("invalid parameter format: %s", pair)
			}
			paramsMap[kv[0]] = kv[1]
		}
	}

	return paramsMap, nil
}
