package image_transform_server_test

import (
	"bytes"
	"errors"
	creatingsymmetry "github.com/Chadius/creating-symmetry"
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
	"reflect"
	"testing"
)

func TestServerUsesPackageSuite(t *testing.T) {
	suite.Run(t, new(ServerUsesPackageSuite))
}

type ServerUsesPackageSuite struct {
	suite.Suite
	request                         *http.Request
	responseRecorder                *httptest.ResponseRecorder
	server                          image_transform_server.TwirpServer
	fakeTransformPackage            *creatingsymmetryfakes.FakeTransformerStrategy
	inputImageData                  []byte
	formulaData                     []byte
	outputSettingsData              []byte
	fakeTransformPackageReturnValue []byte
}

func (suite *ServerUsesPackageSuite) SetupTest() {
	suite.inputImageData = []byte(`images go here`)
	suite.formulaData = []byte(`formula goes here`)
	suite.outputSettingsData = []byte(`outputSettings go here`)
	requestBody := suite.getDataStream()
	suite.request = suite.generateProtobufRequest(requestBody)
	suite.responseRecorder = httptest.NewRecorder()

	suite.fakeTransformPackageReturnValue = []byte(`rules responded`)
	suite.fakeTransformPackage = suite.fakeTransformerStrategyWithResponse(suite.fakeTransformPackageReturnValue)
	suite.server = suite.getServer()
}

func (suite *ServerUsesPackageSuite) fakeTransformerStrategyWithResponse(expectedResponse []byte) *creatingsymmetryfakes.FakeTransformerStrategy {
	fakeTransformerStrategy := creatingsymmetryfakes.FakeTransformerStrategy{}
	fakeTransformerStrategy.ApplyFormulaToTransformImageStub = func(inputImageDataByteStream, formulaDataByteStream, outputSettingsDataByteStream io.Reader, output io.Writer) error {
		output.Write(expectedResponse)
		return nil
	}
	return &fakeTransformerStrategy
}

func (suite *ServerUsesPackageSuite) getServer() image_transform_server.TwirpServer {
	server := transformserver.NewServer(suite.fakeTransformPackage)
	twirpServer := image_transform_server.NewImageTransformerServer(server)
	return twirpServer
}

func (suite *ServerUsesPackageSuite) generateProtobufRequest(requestBody *bytes.Buffer) *http.Request {
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

func (suite *ServerUsesPackageSuite) getDataStream() *bytes.Buffer {
	dataStream := &image_transform_server.DataStreams{
		InputImage:     suite.inputImageData,
		FormulaData:    suite.formulaData,
		OutputSettings: suite.outputSettingsData,
	}

	protobuf, protobufErr := proto.Marshal(dataStream)
	requestBody := bytes.NewBuffer(protobuf)

	require := require.New(suite.T())
	require.Nil(protobufErr)

	return requestBody
}

func (suite *ServerUsesPackageSuite) TestWhenClientMakesRequest_ResponseIsValid() {
	// Act
	suite.server.ServeHTTP(suite.responseRecorder, suite.request)

	// Assert
	response := suite.responseRecorder.Result()

	require := require.New(suite.T())
	require.Equal(200, response.StatusCode, "Status code is wrong")
}

func (suite *ServerUsesPackageSuite) TestWhenClientMakesRequest_PackageIsCalledWithInputData() {
	// Act
	suite.server.ServeHTTP(suite.responseRecorder, suite.request)

	// Assert
	response := suite.responseRecorder.Result()

	require := require.New(suite.T())
	require.Equal(200, response.StatusCode, "Status code is wrong")

	suite.requireFakePackageWasCalledWithExpectedData(require)
}

func (suite *ServerUsesPackageSuite) requireFakePackageWasCalledWithExpectedData(require *require.Assertions) {
	require.Equal(1, suite.fakeTransformPackage.ApplyFormulaToTransformImageCallCount())

	actualInputImageDataByteStream, actualFormulaDataByteStream, actualOutputSettingsDataByteStream, _ := suite.fakeTransformPackage.ApplyFormulaToTransformImageArgsForCall(0)

	actualInputImageData, imageReadErr := ioutil.ReadAll(actualInputImageDataByteStream)
	require.Nil(imageReadErr, "Error while reading input image data from mock object")
	require.Equal(0, bytes.Compare(suite.inputImageData, actualInputImageData), "input image given to mock object is different")

	actualFormulaData, formulaReadErr := ioutil.ReadAll(actualFormulaDataByteStream)
	require.Nil(formulaReadErr, "Error while reading formula data from mock object")
	require.Equal(0, bytes.Compare(suite.formulaData, actualFormulaData), "formula given to mock object is different")

	actualOutputSettingsData, outputSettingsReadErr := ioutil.ReadAll(actualOutputSettingsDataByteStream)
	require.Nil(outputSettingsReadErr, "Error while reading output settings data from mock object")
	require.Equal(0, bytes.Compare(suite.outputSettingsData, actualOutputSettingsData), "output settings given to mock object is different")
}

func (suite *ServerUsesPackageSuite) TestWhenClientMakesRequest_ResponseIsUnmarshalled() {
	// Act
	suite.server.ServeHTTP(suite.responseRecorder, suite.request)

	// Assert
	require := require.New(suite.T())
	suite.requireResponseDataMatches(require)
}

func (suite *ServerUsesPackageSuite) requireResponseDataMatches(require *require.Assertions) {
	output := &image_transform_server.Image{}
	unmarshalErr := proto.Unmarshal(suite.responseRecorder.Body.Bytes(), output)
	require.Nil(unmarshalErr, "Error while unmarshalling response body")
	require.Equal(suite.fakeTransformPackageReturnValue, output.ImageData, "output image received from mock object is different")
}

func (suite *ServerUsesPackageSuite) TestWhenPackageRaisesError_ThenServerReturns500() {
	// Setup
	fakeTransformerStrategy := creatingsymmetryfakes.FakeTransformerStrategy{}
	fakeTransformerStrategy.ApplyFormulaToTransformImageStub = func(inputImageDataByteStream, formulaDataByteStream, outputSettingsDataByteStream io.Reader, output io.Writer) error {
		return errors.New("irrelevant error")
	}
	server := transformserver.NewServer(&fakeTransformerStrategy)
	twirpServer := image_transform_server.NewImageTransformerServer(server)
	responseRecorder := httptest.NewRecorder()

	// Act
	twirpServer.ServeHTTP(responseRecorder, suite.request)

	// Require
	response := responseRecorder.Result()

	require := require.New(suite.T())
	require.Equal(500, response.StatusCode, "Status code is wrong")
}

func (suite *ServerUsesPackageSuite) TestWhenPackagePanics_ThenServerReturns500() {
	// Setup
	var nilObject creatingsymmetry.TransformerStrategy
	fakeTransformerStrategy := creatingsymmetryfakes.FakeTransformerStrategy{}
	fakeTransformerStrategy.ApplyFormulaToTransformImageStub = func(dummyReader1, dummyReader2, dummyReader3 io.Reader, dummyWriter io.Writer) error {
		nilObject.ApplyFormulaToTransformImage(
			dummyReader1,
			dummyReader2,
			dummyReader3,
			dummyWriter,
		)
		return nil
	}
	server := transformserver.NewServer(&fakeTransformerStrategy)
	twirpServer := image_transform_server.NewImageTransformerServer(server)
	responseRecorder := httptest.NewRecorder()

	// Act
	twirpServer.ServeHTTP(responseRecorder, suite.request)

	// Require
	response := responseRecorder.Result()

	require := require.New(suite.T())
	require.Equal(500, response.StatusCode, "Status code is wrong")
}

type InjectTransformerSuite struct {
	suite.Suite
}

func TestInjectTransformerSuite(t *testing.T) {
	suite.Run(t, new(InjectTransformerSuite))
}

func (suite *InjectTransformerSuite) TestDefaultsToProductionImageTransformPackage() {
	// Setup
	productionTransformer := &creatingsymmetry.FileTransformer{}

	// Act
	server := transformserver.NewServer(nil)

	// Assert
	require := require.New(suite.T())
	require.Equal(
		reflect.TypeOf(server.GetTransformer()),
		reflect.TypeOf(productionTransformer),
	)
}

func (suite *InjectTransformerSuite) TestUsesInjectedImageTransformPackage() {
	// Setup
	fakeTransformer := &creatingsymmetryfakes.FakeTransformerStrategy{}

	// Act
	server := transformserver.NewServer(fakeTransformer)

	// Assert
	require := require.New(suite.T())
	require.Equal(
		reflect.TypeOf(server.GetTransformer()),
		reflect.TypeOf(fakeTransformer),
	)
}
