// private/topology/json/json_dualstack_test.go
package json

import (
	"encoding/json"
	"fmt"
	"net"
	"reflect"
	"testing"
)

func TestBRInfoDualStackJSON(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want []string
	}{
		{
			name: "legacy single internal_addr",
			in:   `{"internal_addr":"10.0.0.1:30042","interfaces":{}}`,
			want: []string{"10.0.0.1:30042"},
		},
		{
			name: "new internal_addrs list",
			in: `{"internal_addrs":["[2001:db8::1]:30042","10.0.0.1:30042"],
			       "interfaces":{}}`,
			want: []string{"[2001:db8::1]:30042", "10.0.0.1:30042"},
		},
		{
			name: "both fields: internal_addrs wins",
			in: `{"internal_addr":"10.0.0.1:30042",
			       "internal_addrs":["[2001:db8::1]:30042"],
			       "interfaces":{}}`,
			want: []string{"[2001:db8::1]:30042"},
		},
		{
			name: "neither field",
			in:   `{"interfaces":{}}`,
			want: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var br BRInfo
			if err := json.Unmarshal([]byte(tc.in), &br); err != nil {
				t.Fatalf("unmarshal failed: %v", err)
			}

			got := br.AllInternalAddrs()

			fmt.Printf("\n[%s]\n", tc.name)
			if len(got) == 0 {
				fmt.Println("  No addresses found after unmarshal")
			} else {
				for i, addr := range got {
					host, port, err := net.SplitHostPort(addr)
					if err != nil {
						fmt.Printf("  [%d] raw     : %s  (could not split: %v)\n", i, addr, err)
					} else {
						fmt.Printf("  [%d] address : %s\n", i, host)
						fmt.Printf("  [%d] port    : %s\n", i, port)
					}
				}
			}

			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("AllInternalAddrs()\n  got:  %v\n  want: %v", got, tc.want)
			}
		})
	}
}

