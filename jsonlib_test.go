package jsonlib

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPostJSON(t *testing.T) {

}

func Test_jsonLibrary_GetJSON(t *testing.T) {
	// mock
	type person struct {
		Name string `json:"name"`
	}
	wantValue := &person{Name: "name 123"}
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Send response to be tested
		data, _ := json.Marshal(wantValue)
		rw.Write(data)
	}))
	defer server.Close()
	// run
	type fields struct {
		_client *http.Client
	}
	type args struct {
		url string
		res *person
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantValue *person
		wantErr   bool
	}{
		{
			name: "basic",
			fields: fields{
				_client: server.Client(),
			},
			args: args{
				url: server.URL,
				res: &person{},
			},
			wantValue: wantValue,
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &jsonLibrary{
				_client: tt.fields._client,
			}
			if err := m.GetJSON(tt.args.url, tt.args.res); (err != nil) != tt.wantErr {
				t.Errorf("jsonLibrary.GetJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.args.res.Name != tt.wantValue.Name {
				t.Errorf("jsonLibrary.GetJSON() value = %v, wantValue %v", tt.args.res, tt.wantValue)
			}
		})
	}
}
