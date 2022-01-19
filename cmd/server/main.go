package main

import (
	"github.com/chadius/image-transform-server/internal/transformserver"
	"github.com/chadius/image-transform-server/rpc/transform/github.com/chadius/image_transform_server"
	"net/http"
)

func main() {
	server := transformserver.NewServer(nil)
	twirpHandler := image_transform_server.NewImageTransformerServer(server)

	http.ListenAndServe(":8080", twirpHandler)
}
