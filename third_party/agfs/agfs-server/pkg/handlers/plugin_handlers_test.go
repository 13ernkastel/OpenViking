package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/c4pt0r/agfs/agfs-server/pkg/mountablefs"
	"github.com/c4pt0r/agfs/agfs-server/pkg/plugin/api"
)

func TestManagementEndpointsDisabledByDefault(t *testing.T) {
	mfs := mountablefs.NewMountableFS(api.PoolConfig{})
	handler := NewPluginHandler(mfs, false)
	mux := http.NewServeMux()
	handler.SetupRoutes(mux)

	tests := []struct {
		name   string
		method string
		target string
	}{
		{name: "list mounts", method: http.MethodGet, target: "/api/v1/mounts"},
		{name: "mount plugin", method: http.MethodPost, target: "/api/v1/mount"},
		{name: "unmount plugin", method: http.MethodPost, target: "/api/v1/unmount"},
		{name: "list plugins", method: http.MethodGet, target: "/api/v1/plugins"},
		{name: "load plugin", method: http.MethodPost, target: "/api/v1/plugins/load"},
		{name: "unload plugin", method: http.MethodPost, target: "/api/v1/plugins/unload"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.target, nil)
			rec := httptest.NewRecorder()

			mux.ServeHTTP(rec, req)

			if rec.Code != http.StatusForbidden {
				t.Fatalf("status = %d, want %d", rec.Code, http.StatusForbidden)
			}
		})
	}
}

func TestManagementEndpointsEnabledPreserveExistingValidation(t *testing.T) {
	mfs := mountablefs.NewMountableFS(api.PoolConfig{})
	handler := NewPluginHandler(mfs, true)
	mux := http.NewServeMux()
	handler.SetupRoutes(mux)

	tests := []struct {
		name       string
		method     string
		target     string
		wantStatus int
	}{
		{name: "list mounts", method: http.MethodGet, target: "/api/v1/mounts", wantStatus: http.StatusOK},
		{name: "mount plugin", method: http.MethodPost, target: "/api/v1/mount", wantStatus: http.StatusBadRequest},
		{name: "unmount plugin", method: http.MethodPost, target: "/api/v1/unmount", wantStatus: http.StatusBadRequest},
		{name: "list plugins", method: http.MethodGet, target: "/api/v1/plugins", wantStatus: http.StatusOK},
		{name: "load plugin", method: http.MethodPost, target: "/api/v1/plugins/load", wantStatus: http.StatusBadRequest},
		{name: "unload plugin", method: http.MethodPost, target: "/api/v1/plugins/unload", wantStatus: http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.target, nil)
			rec := httptest.NewRecorder()

			mux.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d", rec.Code, tt.wantStatus)
			}
		})
	}
}
