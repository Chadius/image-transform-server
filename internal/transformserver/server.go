package transformserver

import (
	"bytes"
	"context"
	creatingsymmetry "github.com/Chadius/creating-symmetry"
	"github.com/chadius/image-transform-server/rpc/transform/github.com/chadius/image_transform_server"
)

// Server implements the Transformer service
type Server struct {
	transformer creatingsymmetry.TransformerStrategy
}

// Transform applies the given formula to the image and uses the output settings to return a new image.
func (s *Server) Transform(cts context.Context, data *image_transform_server.DataStreams) (*image_transform_server.Image, error) {

	inputImageDataByteStream := bytes.NewBuffer(data.GetInputImage())
	formulaDataByteStream := bytes.NewBuffer(data.GetFormulaData())
	outputSettingsDataByteStream := bytes.NewBuffer(data.GetOutputSettings())

	var outputImageBuffer bytes.Buffer

	transformErr := s.GetTransformer().ApplyFormulaToTransformImage(inputImageDataByteStream, formulaDataByteStream, outputSettingsDataByteStream, &outputImageBuffer)
	outputImage := &image_transform_server.Image{ImageData: outputImageBuffer.Bytes()}
	return outputImage, transformErr
}

func (s *Server) GetTransformer() creatingsymmetry.TransformerStrategy {
	return s.transformer
}

// NewServer returns a new Server object with the given transformer.
//   Defaults to using the production Transformer if none is given.
func NewServer(transformer creatingsymmetry.TransformerStrategy) *Server {
	var transformerToUse creatingsymmetry.TransformerStrategy
	transformerToUse = &creatingsymmetry.FileTransformer{}
	if transformer != nil {
		transformerToUse = transformer
	}
	return &Server{
		transformer: transformerToUse,
	}
}
