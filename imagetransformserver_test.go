package image_transform_server_test

import (
	"bytes"
	"github.com/chadius/image-transform-server/creatingsymmetryfakes"
	"github.com/chadius/image-transform-server/internal/transformserver"
	"github.com/chadius/image-transform-server/rpc/transform/github.com/chadius/image_transform_server"
	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestExampleTestSuite(t *testing.T) {
	suite.Run(t, new(UsePackageTestSuite))
}

type UsePackageTestSuite struct {
	suite.Suite
}

func (suite *UsePackageTestSuite) TestWhenFilesAreSupplied_ThenCallPackage() {
	// Setup
	imageData := []byte(`images go here`)
	formulaData := []byte(`formula goes here`)
	outputSettingsData := []byte(`outputSettings go here`)
	requestBody := suite.getDataStream(imageData, formulaData, outputSettingsData)
	testRequest := suite.generateProtobufRequest(requestBody)

	expectedResponse := []byte(`rules responded`)
	fakeTransformerStrategy := suite.getFakeTransformerStrategyWithResponse(expectedResponse)
	twirpServer := suite.getTwirpServer(fakeTransformerStrategy)

	responseRecorder := httptest.NewRecorder()

	// Act
	twirpServer.ServeHTTP(responseRecorder, testRequest)

	// Assert
	response := responseRecorder.Result()

	require := require.New(suite.T())
	require.Equal(200, response.StatusCode, "Status code is wrong")

	suite.requireResponseDataMatches(responseRecorder, require, expectedResponse)
	suite.requireFakePackageWasCalledWithExpectedData(require, fakeTransformerStrategy, imageData, formulaData, outputSettingsData)
}

func (suite *UsePackageTestSuite) requireFakePackageWasCalledWithExpectedData(require *require.Assertions, fakeTransformerStrategy *creatingsymmetryfakes.FakeTransformerStrategy, imageData, formulaData, outputSettingsData []byte) {
	require.Equal(1, fakeTransformerStrategy.ApplyFormulaToTransformImageCallCount())

	actualInputImageDataByteStream, actualFormulaDataByteStream, actualOutputSettingsDataByteStream, _ := fakeTransformerStrategy.ApplyFormulaToTransformImageArgsForCall(0)

	actualInputImageData, imageReadErr := ioutil.ReadAll(actualInputImageDataByteStream)
	require.Nil(imageReadErr, "Error while reading input image data from mock object")
	require.Equal(0, bytes.Compare(imageData, actualInputImageData), "input image given to mock object is different")

	actualFormulaData, formulaReadErr := ioutil.ReadAll(actualFormulaDataByteStream)
	require.Nil(formulaReadErr, "Error while reading formula data from mock object")
	require.Equal(0, bytes.Compare(formulaData, actualFormulaData), "formula given to mock object is different")

	actualOutputSettingsData, outputSettingsReadErr := ioutil.ReadAll(actualOutputSettingsDataByteStream)
	require.Nil(outputSettingsReadErr, "Error while reading output settings data from mock object")
	require.Equal(0, bytes.Compare(outputSettingsData, actualOutputSettingsData), "output settings given to mock object is different")
}

func (suite *UsePackageTestSuite) requireResponseDataMatches(responseRecorder *httptest.ResponseRecorder, require *require.Assertions, expectedResponse []byte) {
	output := &image_transform_server.Image{}
	unmarshalErr := proto.Unmarshal(responseRecorder.Body.Bytes(), output)
	require.Nil(unmarshalErr, "Error while unmarshalling response body")
	require.Equal(expectedResponse, output.ImageData, "output image received from mock object is different")
}

func (suite *UsePackageTestSuite) getFakeTransformerStrategyWithResponse(expectedResponse []byte) *creatingsymmetryfakes.FakeTransformerStrategy {
	fakeTransformerStrategy := creatingsymmetryfakes.FakeTransformerStrategy{}
	fakeTransformerStrategy.ApplyFormulaToTransformImageStub = func(inputImageDataByteStream, formulaDataByteStream, outputSettingsDataByteStream io.Reader, output io.Writer) error {
		output.Write(expectedResponse)
		return nil
	}
	return &fakeTransformerStrategy
}

func (suite *UsePackageTestSuite) getTwirpServer(fakeTransformerStrategy *creatingsymmetryfakes.FakeTransformerStrategy) image_transform_server.TwirpServer {
	server := &transformserver.Server{
		Transformer: fakeTransformerStrategy,
	}
	twirpServer := image_transform_server.NewImageTransformerServer(server)
	return twirpServer
}

func (suite *UsePackageTestSuite) generateProtobufRequest(requestBody *bytes.Buffer) *http.Request {
	testRequest, newRequestErr := http.NewRequest(
		http.MethodPost,
		"/twirp/chadius.imageTransformServer.ImageTransformer/Transform",
		requestBody,
	)
	require := require.New(suite.T())
	require.Nil(newRequestErr)
	testRequest.Header.Set("Content-Type", "application/protobuf")
	return testRequest
}

func (suite *UsePackageTestSuite) getDataStream(imageData []byte, formulaData []byte, outputSettingsData []byte) *bytes.Buffer {
	dataStream := &image_transform_server.DataStreams{
		InputImage:     imageData,
		FormulaData:    formulaData,
		OutputSettings: outputSettingsData,
	}

	protobuf, protobufErr := proto.Marshal(dataStream)
	requestBody := bytes.NewBuffer(protobuf)

	require := require.New(suite.T())
	require.Nil(protobufErr)

	return requestBody
}
