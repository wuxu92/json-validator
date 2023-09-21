package validator

import (
	"testing"
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

func TestValidate(t *testing.T) {
	args := []struct {
		name    string
		data    []byte
		wantErr bool
		opts    []Option
	}{
		{
			"not_null tag",
			[]byte(`{"id": null, "not_null_id": null, "not_zero_id": 1, "required_id": 0, "embed": null}`),
			true,
			nil,
		},
		{
			"required tag",
			[]byte(`{"id": null, not_null_id": 0, "not_zero_id": 1, "embed": null}`),
			true,
			nil,
		},
		{
			"redundant option",
			[]byte(`{"id": null, "foo": "", "not_null_id": 0, "not_zero_id": 1, "required_id": 0, "embed": null}`),
			true,
			[]Option{WithNoRedundant},
		},
		{
			"empty embed field required",
			[]byte(`{"id": null, "not_null_id": 0, "not_zero_id": 1, "required_id": 0, "embed": {}}`),
			true,
			nil,
		},
		{
			"embed not_null tag",
			[]byte(`{"id": null, "not_null_id": 0, "not_zero_id": 1, "required_id": 0, "embed": {"not_null_id": null, "not_zero_id": 1, "required_id": 0}}`),
			true,
			nil,
		},
		{
			"embed required tag",
			[]byte(`{"id": null, "not_null_id": 0, "not_zero_id": 1, "required_id": 0, "embed": {"not_null_id": 0, "not_zero_id": 1}}`),
			true,
			nil,
		},
		{
			"embed redundant",
			[]byte(`{"id": null, "not_null_id": 0, "not_zero_id": 1, "required_id": 0, "embed": {"not_null_id": 0, "foo": "bar", "not_zero_id": 1, "required_id": 0}}`),
			true,
			[]Option{WithNoRedundant},
		},
		{
			"shall pass",
			[]byte(`{"id": null, "not_null_id": 0, "not_zero_id": 1, "required_id": 0, "embed": {"not_null_id": 0, "not_zero_id": 1, "required_id": 0}}`),
			false,
			nil,
		},
	}
	for idx, arg := range args {
		err := Validate(arg.data, TestStruct{}, arg.opts...)
		t.Logf("%d(%s) test with err: %+v", idx, arg.name, err)

		if (err != nil) != arg.wantErr {
			t.Fatalf("[FAILED] validata %d(%s) with want err %t, but got %t(%v)", idx, arg.name, arg.wantErr, err != nil, err)
		}
	}

}
