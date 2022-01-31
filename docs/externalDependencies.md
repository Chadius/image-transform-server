# Use counterfeiter to get test objects

```bash
go run github.com/maxbrunsfeld/counterfeiter/v6 -o ./creatingsymmetryfakes github.com/Chadius/creating-symmetry.TransformerStrategy
```

# Protobuf and Twirp
Add the twirp files to tools.go:
```go
	_ "github.com/twitchtv/twirp/protoc-gen-twirp"
	_ "google.golang.org/protobuf/cmd/protoc-gen-go"
```

Set up your paths:
```bash
export GOBIN=$(PWD)/_bin
export PATH=$GOBIN:$PATH
```

Call `go install` to install both tools. This will support twirp and go protobuf files.

Write your `service.proto` file. Dashes don't really work so use all lowercase.

When it's time, run this to compile the protobuf and twirp files.
```bash
protoc -I=./ --twirp_out=./rpc --go_out=./rpc service.proto
```

Now I can look in `./rpc/transform/github.com/chadius/image_transform_server` for the `service.twirp.go` file.

# Testing

The test framework acts as a client, connecting to your server implementation.

Make the server and a handler to that server.
```go
    server := &transformserver.Server{
        Transformer: &fakeTransformerStrategy,
    }
    twirpServer := image_transform_server.NewImageTransformerServer(server)
```

Marshal data into a protobuf object
```go
	dataStream := &image_transform_server.DataStreams{
		InputImage:     imageData,
		FormulaData:    formulaData,
		OutputSettings: outputSettingsData,
	}

	protobuf, protobufErr := proto.Marshal(dataStream)
	requestBody := bytes.NewBuffer(protobuf)
```

Make an HttpRequest object. It takes
- a HTTP method (going to be POST for these RPC callers)
- the route you want to test (Expected format: "[<prefix>]/<package>.<Service>/<Method>")
  - prefix is usually `/twirp`
  - package, service and Method were defined in your `.proto` file
- an io.Reader object that holds the body (nil is fine if you have nothing to post)
- You also have to set the header so the server knows if it' JSON or protobuf.

```go
	testRequest, newRequestErr := http.NewRequest(
		http.MethodPost,
        "/twirp/chadius.imageTransformServer.ImageTransformer/Transform",
		requestBody,
	)
    testRequest.Header.Set("Content-Type", "application/protobuf")
```

Now you can make a `httptest.NewRecorder()` object to hold the results.
``` go
	responseRecorder := httptest.NewRecorder()
```

Now make the server request using ServeHTTP:
``` go
	// Act
	twirpServer.ServeHTTP(responseRecorder, testRequest)
```

You can Assert against the responseRecorder.
``` go
	// Assert
    response := responseRecorder.Result()
	assert.Equal(200, response.StatusCode, "Status code is wrong")
```

You have to unmarshal the response body.
```go
    output := &image_transform_server.Image{}
	unmarshalErr := proto.Unmarshal(responseRecorder.Body.Bytes(), output)
	require.Nil(unmarshalErr, "Error while unmarshalling response body")
	require.Equal(expectedResponse, output.ImageData, "output image received from mock object is different")
```