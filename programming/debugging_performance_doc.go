// Go wiki for debugging performance issue:
// https://go.dev/wiki/Performance#memory-allocator-trace
//
// [before diving into it]
// 1. use go tool to identify various types of hotspots(CPU, IO, memory)
// 2. help developer identify obvious performance defects in the program though profiles
// 3. there are bounding factors for performance, no much you can do for the performance
//
// [1. cpu profiler] cpu分析器
// there are mainly three ways of doing it:
//
// [1.1 go test]
// a. go test -run=none -bench=YourBenchFunction -cpuprofile=cpu.profile .
// b. go tool pproof --text test-binary-file cpu.profile
//
// [1.2 http server]
//
//	a. import _ "net/http/pproof"
//	b. go tool pproof --text server-banary-file http://{host}:{port}/debug/pproof/profile
//
// [1.3 manually generate profile]
//
//	a. import runtime/pproof
//	b. if *flagCpuprofile != "" {
//	   f, err := os.Create(*flagCpuprofile)
//	   if err != nil {
//	   log.Fatal(err)
//	   }
//	   pprof.StartCPUProfile(f)
//	   defer pprof.StopCPUProfile()
//	   }
//
// [2. Memory Profiler]
//
//	a. import runtime/pproof
//	b. if *flagCpuprofile != "" {
//	   f, err := os.Create(*flagCpuprofile)
//	   if err != nil {
//	   log.Fatal(err)
//	   }
//	   pprof.StartCPUProfile(f)
//	   defer pprof.StopCPUProfile()
//	   }
package main
