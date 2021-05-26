# Go source code comments unmarshaler



```go
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
```
