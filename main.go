package main

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/activeterm"
	"github.com/charmbracelet/wish/bubbletea"
	"github.com/charmbracelet/wish/logging"
	"github.com/trafalgar-2006/ssh-portfolio/config"
	"github.com/trafalgar-2006/ssh-portfolio/views"
)

//go:embed index.html
var indexHTML embed.FS

func main() {
	// Load content.yaml — falls back to hardcoded data if file not found
	if err := config.Load("content.yaml"); err != nil {
		log.Printf("content.yaml not found, using hardcoded content: %v", err)
	} else {
		views.LoadFromConfig()
	}

	sshEnabled := os.Getenv("SSH_ENABLED")
	if sshEnabled == "" {
		sshEnabled = "false"
	}

	if sshEnabled == "true" {
		runSSHServer()
	} else {
		runLocalTUI()
	}
}

func runLocalTUI() {
	p := tea.NewProgram(
		NewModel(nil),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running TUI: %v\n", err)
		os.Exit(1)
	}
}

func runSSHServer() {
	host := os.Getenv("HOST")
	if host == "" {
		host = "0.0.0.0"
	}
	port := os.Getenv("SSH_PORT")
	if port == "" {
		port = "23234"
	}

	s, err := wish.NewServer(
		wish.WithAddress(net.JoinHostPort(host, port)),
		wish.WithHostKeyPath(".ssh/id_ed25519"),
		// Accept ALL connections — public portfolio, no auth needed
		wish.WithPublicKeyAuth(func(_ ssh.Context, _ ssh.PublicKey) bool {
			return true
		}),
		wish.WithMiddleware(
			bubbletea.Middleware(teaHandler),
			activeterm.Middleware(),
			logging.Middleware(),
		),
	)
	if err != nil {
		log.Fatalf("Could not create SSH server: %v", err)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// HTTP server: web portfolio at / and health check at /health
	go func() {
		mux := http.NewServeMux()

		// Serve the web portfolio (index.html embedded in the binary)
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			data, err := indexHTML.ReadFile("index.html")
			if err != nil {
				http.Error(w, "portfolio not found", http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.Header().Set("Cache-Control", "public, max-age=3600")
			w.Write(data)
		})

		// Health check for UptimeRobot / Railway
		mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		})

		log.Println("Web portfolio + health check listening on :8080")
		http.ListenAndServe(":8080", mux)
	}()

	log.Printf("Starting SSH server on %s:%s", host, port)
	log.Printf("Connect with: ssh localhost -p %s", port)

	go func() {
		if err := s.ListenAndServe(); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
			log.Fatalf("SSH server error: %v", err)
		}
	}()

	<-done
	log.Println("Shutting down SSH server...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		log.Fatalf("Could not shutdown server: %v", err)
	}
}

func teaHandler(s ssh.Session) (tea.Model, []tea.ProgramOption) {
	renderer := bubbletea.MakeRenderer(s)
	m := NewModel(renderer)
	return m, []tea.ProgramOption{
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	}
}
