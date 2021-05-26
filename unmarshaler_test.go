package go_comments_unmarshaler

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

func TestUnmarshalPackage(t *testing.T) {
	type BasicExample struct {
		Fetcher         string `comment:"Fetcher"`
		FetchOrders     string `comment:"Fetcher.FetchOrders"`
		FetchEmails     string `comment:"Fetcher.FetchEmails"`
		FetchUsers      string `comment:"Fetcher.fetchUsers"`
		FetchNoComments string `comment:"Fetcher.FetchNoComments"`
		F1              string `comment:"F1"`
		F2              string `comment:"F2"`
	}

	type args struct {
		pathToModule string
		result       BasicExample
	}
	tests := []struct {
		name       string
		args       args
		wantErr    bool
		wantResult interface{}
	}{
		{
			name: "basic example", args: args{
				pathToModule: "./testdata",
				result:       BasicExample{},
			},
			wantResult: BasicExample{
				Fetcher:         "Fetcher is a main fetcher for this module\n",
				FetchOrders:     "FetchOrders fetching orders for me.\nThere is also second line\n",
				FetchEmails:     "FetchEmails is for fetching emails.\n\nThis file is placed in second file.\n",
				FetchUsers:      "fetchUsers is private function.\n",
				FetchNoComments: "",
				F1:              "F1 comment\n",
				F2:              "F2 comment\n",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := UnmarshalPackage(tt.args.pathToModule, &tt.args.result); (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalPackage() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(tt.args.result, tt.wantResult) {
				t.Errorf("UnmarshalPackage() error = %+v, wantResult %+v", tt.args.result, tt.wantResult)
				//var b bytes.Buffer
				//err := diff.Text("result", "want",
				//	fmt.Sprintf("%+v", tt.args.result),
				//	fmt.Sprintf("%+v", tt.wantResult),
				//	&b,
				//)
				//if err != nil {
				//	t.Error(err)
				//}
				//t.Log(b.String())
			}
		})
	}
}

func TestUnmarshalModule(t *testing.T) {
	type Example struct {
		Fetcher         string `comment:"Fetcher"`
		FetchOrders     string `comment:"Fetcher.FetchOrders"`
		FetchEmails     string `comment:"Fetcher.FetchEmails"`
		FetchUsers      string `comment:"Fetcher.fetchUsers"`
		FetchNoComments string `comment:"Fetcher.FetchNoComments"`
		Module1         struct {
			PublicFunc  string `comment:"Module1Func"`
			PrivateFunc string `comment:"module1PrivateFunc"`
			SomeType    string `comment:"module2/SomeType"`
			Module2     struct {
				Interface       string `comment:"Interface"`
				Implementations map[string]struct {
					Implementation string `comment:"Implementation"`
					ImplementsFunc string `comment:"Implementation.Implements"`
				} `comment:"*"`
			} `comment:"module2"`
		} `comment:"module1"`
	}

	type args struct {
		pathToModule string
		result       Example
	}
	tests := []struct {
		name       string
		args       args
		wantErr    bool
		wantResult interface{}
	}{
		{
			name: "basic example", args: args{
				pathToModule: "./testdata",
				result:       Example{},
			},
			wantResult: Example{
				Fetcher:         "Fetcher is a main fetcher for this module\n",
				FetchOrders:     "FetchOrders fetching orders for me.\nThere is also second line\n",
				FetchEmails:     "FetchEmails is for fetching emails.\n\nThis file is placed in second file.\n",
				FetchUsers:      "fetchUsers is private function.\n",
				FetchNoComments: "",
				Module1: struct {
					PublicFunc  string `comment:"Module1Func"`
					PrivateFunc string `comment:"module1PrivateFunc"`
					SomeType    string `comment:"module2/SomeType"`
					Module2     struct {
						Interface       string `comment:"Interface"`
						Implementations map[string]struct {
							Implementation string `comment:"Implementation"`
							ImplementsFunc string `comment:"Implementation.Implements"`
						} `comment:"*"`
					} `comment:"module2"`
				}{
					SomeType:    "SomeType is some type\n",
					PublicFunc:  "Module1Func is a module one function\n",
					PrivateFunc: "module1PrivateFunc is module private func.\n",
					Module2: struct {
						Interface       string `comment:"Interface"`
						Implementations map[string]struct {
							Implementation string `comment:"Implementation"`
							ImplementsFunc string `comment:"Implementation.Implements"`
						} `comment:"*"`
					}{
						Interface: "Interface is a main interface\n",
						Implementations: map[string]struct {
							Implementation string `comment:"Implementation"`
							ImplementsFunc string `comment:"Implementation.Implements"`
						}{"module3": {
							Implementation: "Implementation is an implementation in module 3\n",
							ImplementsFunc: "Implements is an implementation of Implementation in module 3\n",
						}, "module4": {
							Implementation: "Implementation is an implementation in module 4\n",
							ImplementsFunc: "Implements is an implementation of Implementation in module 4\n",
						}},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := UnmarshalModule(tt.args.pathToModule, &tt.args.result); (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalModule() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(tt.args.result, tt.wantResult) {
				t.Errorf("UnmarshalModule() error = %+v, wantResult %+v", tt.args.result, tt.wantResult)
				//var b bytes.Buffer
				//err := diff.Text("result", "want",
				//	fmt.Sprintf("%+v", tt.args.result),
				//	fmt.Sprintf("%+v", tt.wantResult),
				//	&b,
				//)
				//if err != nil {
				//	t.Error(err)
				//}
				//t.Log(b.String())
			}
		})
	}
}

func ExampleUnmarshalModule() {
	type ModuleDoc struct {
		// Only string fields is supported
		ClientDoc string `comment:"Client"`
		// For fetching docs about methods, use Receiver.MethodName tag.
		ClientDo string `comment:"Client.Do"`
	}
	type Docs struct {
		// only map[string]struct is supported
		Modules map[string]ModuleDoc `comment:"*"`
	}

	var result Docs
	err := UnmarshalModule("./testdata", &result)
	if err != nil {
		panic(err)
	}
	v, _ := json.MarshalIndent(result, "", " ")
	fmt.Println(string(v))

	// Output:
	// {
	//  "Modules": {
	//   "client": {
	//    "ClientDoc": "Client is client from package `client`\n",
	//    "ClientDo": "Do some stuff\n"
	//   },
	//   "module1": {
	//    "ClientDoc": "And this Client is from module1\n",
	//    "ClientDo": "Comment for another Do\n"
	//   }
	//  }
	// }
}
