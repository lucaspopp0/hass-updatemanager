package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/danielgtaylor/huma/v2/humacli"
	"github.com/go-chi/chi/v5"
	"github.com/spf13/cobra"

	"github.com/lucaspopp0/hass-update-manager/update-manager/api/middleware"
	"github.com/lucaspopp0/hass-update-manager/update-manager/ble"
	"github.com/lucaspopp0/hass-update-manager/update-manager/config"
	"github.com/lucaspopp0/hass-update-manager/update-manager/homeassistant"
	"github.com/lucaspopp0/hass-update-manager/update-manager/model"
	"github.com/lucaspopp0/hass-update-manager/update-manager/util"
)

const (
	envSupervisorToken = "SUPERVISOR_TOKEN"
	envLocal           = "LOCAL"
	envPort            = "PORT"

	defaultServerPort = "8124"
)

type server struct {
	humacli.CLI
	router *chi.Mux
	ha     homeassistant.API

	mCfg sync.Mutex
	cfg  config.Config

	mExecutables sync.Mutex
	executables  homeassistant.Executables

	bleService *ble.Service

	scripts []string
}

func (s *server) onStart() {
	fmt.Println("Starting up...")

	cfg, err := config.FromFile()
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err.Error())
	} else {
		s.cfg = *cfg
	}

	if cfg.Switches == nil {
		cfg.Switches = map[string]model.Switch{}

		err = cfg.WriteFile()
		if err != nil {
			fmt.Printf("Failed to init config: %v\n", err.Error())
		}
	}

	s.executables, err = s.ha.ListExecutables()

	if err != nil {
		fmt.Printf("Home assistant service call failed: %v\n", err.Error())
	}

	fmt.Println("Initializing BLE service...")
	s.bleService, err = ble.NewService()
	if err != nil {
		fmt.Printf("Failed to initialize BLE service: %v\n", err.Error())
	}

	addonInfo, err := s.ha.GetAddOnInfo("self")
	if err != nil {
		fmt.Printf("Failed to fetch add-on info: %v\n", err.Error())
	} else {
		fmt.Printf("Add-on info: %v\n", addonInfo)
	}

	port := util.GetEnv(envPort, defaultServerPort)

	fmt.Printf("Starting server on port %v...\n", port)
	http.ListenAndServe(fmt.Sprintf(":%s", port), s.router)
}

func NewServer() humacli.CLI {
	supervisorToken := os.Getenv(envSupervisorToken)

	s := &server{}

	var api huma.API

	s.CLI = humacli.New(func(hooks humacli.Hooks, o *struct{}) {
		s.ha = homeassistant.NewAPI(homeassistant.APIConfig{
			SupervisorToken: supervisorToken,
		})

		// Create a new router & API.
		s.router = chi.NewMux()

		s.router.Use(AllowCORS)
		s.router.Use(SiteMiddleware(os.Getenv(envLocal) == "true"))

		cfg := huma.DefaultConfig("Smart Switches", "")
		cfg.DocsPath = "/api/docs"
		cfg.SchemasPath = "/api/schemas"
		cfg.OpenAPIPath = "/api/openapi"

		cfg.Servers = []*huma.Server{
			{
				URL:         "",
				Description: "Site",
			},
		}

		api = humachi.New(
			s.router,
			cfg,
		)

		api.UseMiddleware(middleware.Logger)

		huma.AutoRegister(api, s)

		hooks.OnStart(s.onStart)
	})

	// Add a command to print out the OpenAPI doc without starting the server.
	s.CLI.Root().AddCommand(&cobra.Command{
		Use:   "get-openapi",
		Short: "Dump OpenAPI to stdout",
		Long:  "Generate the OpenAPI spec and dump it to stdout without starting the server. Redirect the output to save to a file.",
		Run: func(cmd *cobra.Command, args []string) {
			b, _ := json.MarshalIndent(api.OpenAPI(), "", "  ")
			fmt.Println(string(b))
		}},
	)

	return s
}
