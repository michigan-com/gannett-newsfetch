package gannettApi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
)

func dumpJSONFromReader(intro string, prefix string, body io.Reader) io.Reader {
	original, err := ioutil.ReadAll(body)
	if err != nil {
		panic(err)
	}

	var prettyBuf bytes.Buffer
	json.Indent(&prettyBuf, original, prefix, "  ")
	pretty := strings.TrimSpace(prettyBuf.String())

	fmt.Printf("%s\n%s\n", intro, pretty)

	return bytes.NewBuffer(original)
}
