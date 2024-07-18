package main

import "gonetworker/core"

func main() {
	core.Send("./example/docs/gonetwork.config.json")
	//if len(os.Args) != 3 {
	//	fmt.Println("Usage: ./main [-g|-r] <file>")
	//	fmt.Println("  -g: Generate config file from Go source file")
	//	fmt.Println("  -r: Run requests using existing config file")
	//	os.Exit(1)
	//}
	//
	//switch os.Args[1] {
	//case "-g":
	//	core.GenerateConfig(os.Args[2])
	//case "-r":
	//	core.RequestTo()
	//default:
	//	fmt.Println("Invalid option. Use -g to generate config or -r to run requests.")
	//	os.Exit(1)
	//}
}
