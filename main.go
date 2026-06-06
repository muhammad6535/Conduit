package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/muhammad6535/conduit/internal/config"
	"github.com/muhammad6535/conduit/internal/cost"
	"github.com/muhammad6535/conduit/internal/provider"
	"github.com/muhammad6535/conduit/internal/router"
	"github.com/muhammad6535/conduit/internal/server"
)

var (
	version = "0.1.0"
	commit  = "unknown"
	date    = "unknown"
)

func main() {
	configPath := flag.String("config", "config.yaml", "Path to config file")
	showVersion := flag.Bool("version", false, "Show version")
	initConfig := flag.Bool("init", false, "Generate default config")
	flag.Parse()

	if *showVersion {
		printVersion()
		return
	}

	if *initConfig {
		generateConfig()
		return
	}

	printBanner()

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf(" Failed to load config: %v", err)
	}

	providers := make(map[string]provider.Provider)
	for _, pc := range cfg.Providers {
		var p provider.Provider
		switch pc.Type {
		case "openai":
			p = provider.NewOpenAI(pc.Name, pc.APIKey, pc.BaseURL)
		case "anthropic":
			p = provider.NewAnthropic(pc.Name, pc.APIKey, pc.BaseURL)
		case "mock":
			p = provider.NewMock(pc.Name)
		default:
			log.Fatalf(" Unknown provider type: %s", pc.Type)
		}
		providers[pc.Name] = p
		fmt.Printf(" • Registered: %s (%s)\n", pc.Name, pc.Type)
	}

	rtr, err := router.New(providers, &cfg.Routing)
	if err != nil {
		log.Fatalf(" Failed to create router: %v", err)
	}

	tracker := cost.NewTracker()
	h := server.NewHandler(rtr, tracker)

	mux := chi.NewRouter()
	mux.Use(chimw.Logger)
	mux.Use(chimw.Recoverer)
	mux.Use(chimw.RealIP)

	mux.Get("/health", h.Health)
	mux.Get("/v1/models", h.ListModels)
	mux.Post("/v1/chat/completions", h.ChatCompletions)
	mux.Get("/v1/stats", h.Stats)

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 180 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		fmt.Printf("\n")
		fmt.Printf(" • Listening: %s\n", addr)
		fmt.Printf(" • Default:   %s → %s\n", cfg.Routing.DefaultProvider, cfg.Routing.DefaultModel)
		fmt.Printf(" • Providers: %d configured\n", len(cfg.Providers))
		fmt.Printf("\n")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf(" Server error: %v", err)
		}
	}()

	<-ctx.Done()
	fmt.Printf("\n Shutting down gracefully...\n")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf(" Shutdown error: %v", err)
	}

	if monthly := tracker.MonthlyCost(); monthly > 0 {
		fmt.Printf("\n === Usage Summary ===\n")
		fmt.Printf(" Total cost this session: $%.4f\n", monthly)
		fmt.Printf(" Total requests:          %d\n", len(tracker.Records()))
	}
	fmt.Printf(" Conduit stopped.\n")
}

func printBanner() {
	fmt.Printf("\n")
	fmt.Printf("  ╔═══════════════════════════════════════════════╗\n")
	fmt.Printf("  ║          Conduit AI Gateway v%-5s            ║\n", version)
	fmt.Printf("  ║   The AI Gateway That Pays for Itself        ║\n")
	fmt.Printf("  ╚═══════════════════════════════════════════════╝\n")
}

func printVersion() {
	info, ok := debug.ReadBuildInfo()
	if ok {
		fmt.Printf("Conduit v%s\n", version)
		fmt.Printf("Module: %s\n", info.Main.Path)
	}
	fmt.Printf("Commit: %s\n", commit)
	fmt.Printf("Built:  %s\n", date)
}

func generateConfig() {
	yaml := `# Conduit AI Gateway Configuration
# Environment variables ($VAR) are automatically expanded.

server:
  host: "0.0.0.0"
  port: 8080

providers:
  - name: "openai"
    type: "openai"
    api_key: "${OPENAI_API_KEY}"
    base_url: "https://api.openai.com/v1"
    models:
      - "gpt-4o"
      - "gpt-4o-mini"
      - "gpt-4-turbo"

  - name: "anthropic"
    type: "anthropic"
    api_key: "${ANTHROPIC_API_KEY}"
    base_url: "https://api.anthropic.com/v1"
    models:
      - "claude-sonnet-4"
      - "claude-haiku-3"

  - name: "local"
    type: "openai"
    api_key: "not-needed"
    base_url: "http://localhost:11434/v1"
    models:
      - "llama3"
      - "mistral"

routing:
  default_provider: "openai"
  default_model: "gpt-4o-mini"
  rules:
    - pattern: "claude-.*"
      provider: "anthropic"
    - pattern: "gpt-.*"
      provider: "openai"
    - pattern: "llama.*|mistral.*"
      provider: "local"

cost:
  max_budget_per_month: 1000.0
  alert_threshold: 80.0
`

	path := "config.yaml"
	if _, err := os.Stat(path); err == nil {
		fmt.Printf(" • %s already exists. Will not overwrite.\n", path)
		return
	}
	if err := os.WriteFile(path, []byte(yaml), 0644); err != nil {
		log.Fatalf(" Failed to write config: %v", err)
	}
	fmt.Printf(" • Created %s\n", path)
	fmt.Printf(" • Edit it with your API keys and run:\n")
	fmt.Printf("   conduit --config %s\n", path)
}

func init() {
	log.SetFlags(0)
}
