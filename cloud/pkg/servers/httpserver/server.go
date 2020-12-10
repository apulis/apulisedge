// Copyright 2020 Apulis Technology Inc. All rights reserved.

package httpserver

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/apulis/ApulisEdge/cloud/pkg/configs"
)

// Start API server for portal
func StartApiServer(config *configs.EdgeCloudConfig) *http.Server {
	port := config.Portal.Http.Port
	router := NewRouter()

	logger.Infof("ApulisEdgeCloud started, listening and serving HTTP on: %d, DebugModel: %t", port, config.DebugModel)
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("listen: %s\n", err)
		}
	}()

	return srv
}

// Stop API server for portal
func StopApiServer(srv *http.Server) {
	logger.Println("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown:", err)
	}

	logger.Println("Server exiting")
}
