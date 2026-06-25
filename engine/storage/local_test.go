package storage

import (
	"context"
	"reflect"
	"testing"

	"github.com/azzz/strata/engine/types"
)

func TestLocalStorage_NewScanner(t *testing.T) {
	type fields struct {
		RootDir string
	}
	type args struct {
		ctx context.Context
		req ScanRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    Scanner
		wantErr bool
	}{
		{
			name: "JSONL scanner creation",
			fields: fields{
				RootDir: "/testdata",
			},
			args: args{
				ctx: context.Background(),
				req: ScanRequest{
					Format: FormatJSONL,
					URI:    "data.jsonl",
					Schema: types.NewSchema(
						types.Column{Name: "node_id", Type: types.KindInt64},
						types.Column{Name: "active", Type: types.KindBool},
						types.Column{Name: "load", Type: types.KindFloat64},
					),
				},
			},
			want: &JSONLScanner{
				path: "/testdata/data.jsonl",
				schema: types.NewSchema(
					types.Column{Name: "node_id", Type: types.KindInt64},
					types.Column{Name: "active", Type: types.KindBool},
					types.Column{Name: "load", Type: types.KindFloat64},
				),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &LocalStorage{
				RootDir: tt.fields.RootDir,
			}
			got, err := s.NewScanner(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("LocalStorage.NewScanner() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LocalStorage.NewScanner() = %v, want %v", got, tt.want)
			}
		})
	}
}
