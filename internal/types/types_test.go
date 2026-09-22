package types

import (
	"encoding/json"
	"testing"
)

func TestFlexStringUnmarshalJSON(t *testing.T) {
	cases := []struct {
		name string
		json string
		want FlexString
	}{
		{"string value", `"5"`, "5"},
		{"number value", `5`, "5"},
		{"float-looking number value", `5.0`, "5.0"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var f FlexString
			if err := json.Unmarshal([]byte(tc.json), &f); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if f != tc.want {
				t.Fatalf("got %q, want %q", f, tc.want)
			}
		})
	}
}

func TestUnconfiguredOnuUnmarshalsNumericOnuTypeID(t *testing.T) {
	data := `{"pon_type":"GPON","board":"1","port":"1","onu":"1","sn":"HWTC1234","onu_type_name":"F660","onu_type_id":11,"olt_id":"1"}`

	var onu UnconfiguredOnu
	if err := json.Unmarshal([]byte(data), &onu); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if onu.OnuTypeID != "11" {
		t.Fatalf("got OnuTypeID %q, want %q", onu.OnuTypeID, "11")
	}
}
