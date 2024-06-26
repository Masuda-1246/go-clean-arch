// Package openapi provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.16.2 DO NOT EDIT.
package openapi

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"net/url"
	"path"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/+xWXW8TRxf+K6t536vIeE1iJLQVF5QbkCoVqeoVjqzJemIvWc8ss+MoVmQpM4sgQNqm",
	"QD9IRatGKAlxiYsiJCgWf+bIJFzxF6oz69iO2RhFymV9Y3vmOR/POc/MmWXii3okOOMqJt4yYUu0HoXM",
	"/p6jlbJktxosVvg3btTrVDaJR76kFedoI0cWadhgxCPLJe44JeKLCisRzykWCrl0pc7imFbtYmnUtkRK",
	"vEVaORJwxSSnYTlmcpHJMpNSyOMhr/UhTgpxUsiJwS9kB8/0MkgjolVW5kKV50WDV7KyuE6rzOFCORYx",
	"iXsxK/xx80HcuOH7LI7HQkmBiwGvOmKRyVDQCfGms+kOIjQ4W4qYr1gmq28Hu58r60xhOivOuIN+3FaO",
	"SBZHgsepnqYLBfzyBVeMW0nRKAoDn6pAcPdmLPi4Akcq83/J5olH/ucO9eoeQd0jHIaM/Rqr05NM0t3Y",
	"PUqMtNCmwmJfBhHmQTzyfnW9d/+P908MJH+BeQPJU0j2wbzBSs4Upk/JIav2E8l8YnA2rC43VK3ggN47",
	"3P3ucKd78Pjth6eboF9DsgvJCphnlm0HklWSIzVGK0za/L9h6twVIRYChv+O+wwqZSUWGHdA7zjU9iBd",
	"gBXjTE31vt/vde6AboPuOKkPB3QH9CPQe87l69cwnd6PXdD7vbv/gP4F9Pbhzove3gbon8CsgXkAxoDe",
	"A71t7XZB/wAr+hgD8/Dg1drB45e9d2vo1fVpGM5Rf8FBbnoP6ZkOmB3byVXsqm6PB1kxU1NktMj9ZoxQ",
	"vMSa+Xw+/8Uoy0t0zq/Y1atKRV/zsElyRDUjtIuVDHgVG9HKkeKplT92904UzCj2bLQyVjT8kTyHpEss",
	"l+IpuZx0r04klWl0Vuw2IOmCeQ167XDrAehnoB+AuQ/6HejfwDxCmhdO3bITp9hEntlWZ0TUvEKiybql",
	"O+jhwLdNe+AAOyVFxKQKUkJ46+P3vJB1qvAscDUzPZQ4pl5lEqvVHweIHtc/DoFbjUCyCvFupD6H+NmB",
	"MzF3k/kqZRHweWFdBcqeQRExfo5GAQ4mJuOU3Pl8ASPjHm55ZCaPSzjGVc3m79KGqrmhqAa2V5FIz9K4",
	"GF6A+dtef/u9u1sH63fAPDz8cw30HWK9S9vxaxXika+sq+ypltWjAc5F0HB6TMYiaHhnTMYWU7/9M/k5",
	"bHFE2JOxCLI6YX5DBqpJvBuzuZEnw2jVUBC0GmN3seBkFg3dGqOhqmGkKrNlP17KKzXmL1y1mOZ/FbUV",
	"/RWStr1vNZhtSBIwnfQ4pDcDFniZNGRIPFJTKvJcNxQ+DWsiVt7FwsWCu3ietGYHzRgXejo0IXkLZhNM",
	"G8wWvmqSVdDtDz9vgn5ip+2OXbyHbx7b3PQ5wGkdz6Ftbyv3yRHSt0Fv4CjVbTCm9/J3nKw4rG9nOvzY",
	"Xc1ka02ef+zeG0asB7FvSS2di5WIwqBas2IKUENL/Q9ptf4NAAD//+zfMePHDAAA",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %w", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	res := make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	resolvePath := PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		pathToFile := url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}
