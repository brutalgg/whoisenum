package output

import (
	"encoding/json"
	"io"
)

func JsonWriteOut(out io.Writer, x interface{}) error {
	o, err := json.MarshalIndent(x, "", "  ")
	if err != nil {
		return err
	}
	_, err = out.Write(o)
	return err
}
