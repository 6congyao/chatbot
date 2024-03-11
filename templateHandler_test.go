package main

import (
	"context"
	"reflect"
	"testing"
)

func TestTemplateHandler_Do(t *testing.T) {
	type fields struct {
		Next Next
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    context.Context
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &TemplateHandler{
				Next: tt.fields.Next,
			}
			got, err := tr.Do(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("TemplateHandler.Do() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TemplateHandler.Do() = %v, want %v", got, tt.want)
			}
		})
	}
}
