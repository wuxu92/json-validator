# A simple json validator for Golang

The `json.Unmarshal` function ignores `null` values which may cause confusing results as discussed [here](https://www.v2ex.com/t/975214). This package provides additional json tag to validate if the json bytes of specific field is `null`, `zero` value or has fields that is not defined in the model.

## Usage

```go
package main

import (
	"fmt"

	"github.com/wuxu92/json-validator"
)

type TestStruct struct {
	ID         int `json:"id"`
	NotNulLID  int `json:"not_null_id,omitempty,not_null"`
	NotZeroID  int `json:"not_zero_id,not_zero"`
	RequiredID int `json:"required_id,required"`
	Embed      struct {
		NotNulLID  int `json:"not_null_id,omitempty,not_null"`
		NotZeroID  int `json:"not_zero_id,not_zero"`
		RequiredID int `json:"required_id,required"`
	} `json:"embed,required"`
}

func main() {
	data := []byte(`{
        "id": 1,
        "not_null_id": null,
        "not_zero_id": 0,
        "required_id": 1,
        "embed": {
            "not_null_id": null,
            "not_zero_id": 0,
            "required_id": 1
        }
    }`)
	if err := Validate(data, TestStruct{}, WithNoRedundant); err != nil {
		fmt.Printf("validate with: %+v", err)
	}
}
```