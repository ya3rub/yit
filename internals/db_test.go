package internals

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestStore_ReadObj(t *testing.T) {
	curr, _ := os.Getwd()
	root_dir := filepath.Dir(curr)
	yit_dir := strings.Join(
		[]string{root_dir, YitMetadataDir},
		string(os.PathSeparator),
	)
	db_dir := strings.Join(
		[]string{yit_dir, "objects"},
		string(os.PathSeparator),
	)
	type fields struct {
		path string
	}
	type args struct {
		objID []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			fields: fields{
				path: db_dir,
			},
			name: "test ob",
			args: args{
				objID: []byte("fd749358c62aa780e8dbaeb08492f36ea6b6dba8"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Store{
				path: tt.fields.path,
			}
			if _, err := s.ReadObj(tt.args.objID); (err != nil) != tt.wantErr {
				t.Errorf(
					"Store.ReadObj() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
			}
		})
	}
}
