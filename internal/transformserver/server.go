package transformserver

import (
	"bytes"
	"context"
	creatingsymmetry "github.com/Chadius/creating-symmetry"
	"github.com/chadius/image-transform-server/rpc/transform/github.com/chadius/image_transform_server"
)

// Server implements the Transformer service
type Server struct {
	Transformer creatingsymmetry.TransformerStrategy
}

// Transform applies the given formula to the image and uses the output settings to return a new image.
func (s *Server) Transform(cts context.Context, data *image_transform_server.DataStreams) (*image_transform_server.Image, error) {

	inputImageDataByteStream := bytes.NewBuffer(data.GetInputImage())
	formulaDataByteStream := bytes.NewBuffer(data.GetFormulaData())
	outputSettingsDataByteStream := bytes.NewBuffer(data.GetOutputSettings())

	var outputImageBuffer bytes.Buffer

	transformErr := s.Transformer.ApplyFormulaToTransformImage(inputImageDataByteStream, formulaDataByteStream, outputSettingsDataByteStream, &outputImageBuffer)
	outputImage := &image_transform_server.Image{ImageData: outputImageBuffer.Bytes()}
	return outputImage, transformErr
}
