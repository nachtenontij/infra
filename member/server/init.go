package server

import (
	"github.com/nachtenontij/infra/base/server"
)

func init() {
	server.RegisterRegisterUrls(RegisterUrls)
	server.RegisterInitializeCollections(InitializeCollections)
	server.RegisterMiddleware(Middleware)
}
