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
	openBrowser            bool
	host                   string
	port                   int
	tileMapServiceName     string
	tileMapServiceParams   string
	telemetryServiceName   string
	telemetryServiceParams string
	radarServiceName       string
	radarServiceParams     string
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the application",
	Run: func(cmd *cobra.Command, args []string) {
		start(openBrowser, host, port, tileMapServiceName, tileMapServiceParams, telemetryServiceName, telemetryServiceParams, radarServiceName, radarServiceParams)
	},
}

func init() {
	rootCmd.AddCommand(startCmd)

	startCmd.Flags().BoolVarP(&openBrowser, "openBrowser", "o", true, "Open browser on start")
	startCmd.Flags().IntVarP(&port, "port", "p", 8080, "Port to bind the server to")
	startCmd.Flags().StringVar(&host, "host", "localhost", "Host to bind the server to")

	startCmd.Flags().StringVar(&tileMapServiceName, "tmservice", "openstreetmap", "tile map service to use (ex fake, openstreetmap, mbtiles)")
	startCmd.Flags().StringVar(&tileMapServiceParams, "tmsparams", "", "tile map service parameters")

	startCmd.Flags().StringVar(&telemetryServiceName, "tservice", "mavlinkserial", "Telemetry service to use (ex fake, mavlinkserial)")
	startCmd.Flags().StringVar(&telemetryServiceParams, "tsparams", "", "Telemetry service parameters")

	startCmd.Flags().StringVar(&radarServiceName, "rservice", "opensky", "Radar service to use (ex fake, opensky)")
	startCmd.Flags().StringVar(&radarServiceParams, "rsparams", "", "Radar service parameters")
}

func start(openBrowser bool, host string, port int, tileMapServiceName, tileMapServiceParams, telemetryServiceName, telemetryServiceParams, radarServiceName, radarServiceParams string) {
	ctx, cancel := context.WithCancel(context.Background())
	serverURL := fmt.Sprintf("http://%s:%d", host, port)

	tileMapServiceParamsMap, err := paramsStringToMap(tileMapServiceParams)
	if err != nil {
		log.Fatalf("Error converting params of tile service: %v", err)
	}
	tileMapService, err := app.CreateTileMapService(tileMapServiceName, tileMapServiceParamsMap)
	if err != nil {
		log.Fatalf("Error creating tile service: %v", err)
	}

	telemetryServiceParamsMap, err := paramsStringToMap(telemetryServiceParams)
	if err != nil {
		log.Fatalf("Error converting telemetry params: %v", err)
	}
	telemetryService, err := app.CreateTelemetryService(telemetryServiceName, telemetryServiceParamsMap)
	if err != nil {
		log.Fatalf("Error creating telemetry service: %v", err)
	}
	telemetryService.Initialize(ctx)

	radarServiceParamsMap, err := paramsStringToMap(radarServiceParams)
	if err != nil {
		log.Fatalf("Error converting radar params: %v", err)
	}
	radarService, err := app.CreateRadarService(radarServiceName, radarServiceParamsMap)
	if err != nil {
		log.Fatalf("Error creating radar service: %v", err)
	}

	handleRequests := v1.HandleRequests{
		TelemetryService: telemetryService,
		RadarService:     radarService,
		TileMapService:   tileMapService,
	}

	http.HandleFunc("/api/v1/session", handleRequests.HandleSession)
	http.HandleFunc("/api/v1/radar", handleRequests.HandleRadar)
	http.HandleFunc("/api/v1/telemetry", handleRequests.HandleTelemetry)
	http.HandleFunc("/api/v1/tile/", handleRequests.HandleTileMap)

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

	if openBrowser {
		err = openUrlOnBrowser(serverURL)
		if err != nil {
			log.Printf("Error opening browser: %v", err)
		}
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
