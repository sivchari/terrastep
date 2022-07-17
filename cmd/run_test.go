package cmd

import (
	"testing"

	_ "embed"

	"github.com/google/go-cmp/cmp"
)

func TestParseYML(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		want    *Config
		wantErr bool
	}{
		{
			name: "parse yml",
			path: "testdata/task.yml",
			want: &Config{
				Tasks: []*Task{
					{
						Name: "run a",
						Tactics: []string{
							"validate",
							"fmt",
							"plan",
						},
						Steps: []string{
							"a",
							"b",
						},
					},
					{
						Name: "run c",
						Tactics: []string{
							"validate",
							"fmt",
							"plan",
						},
						Steps: []string{
							"c",
							"d",
						},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := parseyml(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseyml() err = %s\n", err.Error())
				return
			}
			if diff := cmp.Diff(tt.want, cfg); diff != "" {
				t.Errorf("mismatch (-want +got):\n%s", diff)
				return
			}
		})
	}
}
