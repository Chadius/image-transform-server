package main

import (
	"github.com/chadius/image-transform-server/internal/transformserver"
	"github.com/chadius/image-transform-server/rpc/transform/github.com/chadius/image_transform_server"
	"github.com/rs/cors"
	"net/http"
)

func main() {
	server := transformserver.NewServer(nil)
	twirpServer := image_transform_server.NewImageTransformerServer(server)
	handler := cors.Default().Handler(twirpServer)
	http.ListenAndServe(":8080", handler)
}
