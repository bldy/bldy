// Copyright 2015 Sevki <s@sevki.org>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pretty // import "sevki.org/x/pretty"

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// JSON prints a struct neatly in JSON format
func JSON(x interface{}) string {
	b, err := json.Marshal(x)
	if err != nil {
		return fmt.Sprintf("JSON parse error: %s", err)
	}

	var prettyJSON bytes.Buffer
	error := json.Indent(&prettyJSON, b, "", "\t")
	if error != nil {

		return fmt.Sprintf("JSON parse error: %s", err)
	}
	return prettyJSON.String()
}
