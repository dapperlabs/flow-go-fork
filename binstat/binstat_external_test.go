package binstat_test

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"runtime/pprof"
	"sync"
	"testing"
	"time"

	"github.com/onflow/flow-go/binstat"
	"github.com/stretchr/testify/require"
)

// go fmt binstat/*.go ; pushd binstat ; GO111MODULE=on go test -v -coverprofile=coverage.txt -covermode=atomic --tags relic ./... ; go tool cover -func=coverage.txt ; popd

/*
 * NOTE: The code below is inspired by the goroutine.go here [1] [2].
 * [1] https://eng.uber.com/pprof-go-profiler/
 * [2] https://github.com/chabbimilind/GoPprofDemo/blob/master/goroutine.go
 */

const loops = 2
const tries = 3
const funcs = 6
const mechs = 2

var wg sync.WaitGroup
var el [loops][tries][mechs][funcs]string

// each function f1-f6 runs the same function f and adds its wall-clock execution time to a table of elapsed times
func f1(outerFuncName string, f func(string) time.Duration, loop int, try int, i int) {
	defer wg.Done()
	el[loop][try][0][i] = fmt.Sprintf("%.02f", f(outerFuncName).Seconds())
}
func f2(outerFuncName string, f func(string) time.Duration, loop int, try int, i int) {
	defer wg.Done()
	el[loop][try][0][i] = fmt.Sprintf("%.02f", f(outerFuncName).Seconds())
}
func f3(outerFuncName string, f func(string) time.Duration, loop int, try int, i int) {
	defer wg.Done()
	el[loop][try][0][i] = fmt.Sprintf("%.02f", f(outerFuncName).Seconds())
}
func f4(outerFuncName string, f func(string) time.Duration, loop int, try int, i int) {
	defer wg.Done()
	el[loop][try][0][i] = fmt.Sprintf("%.02f", f(outerFuncName).Seconds())
}
func f5(outerFuncName string, f func(string) time.Duration, loop int, try int, i int) {
	defer wg.Done()
	el[loop][try][0][i] = fmt.Sprintf("%.02f", f(outerFuncName).Seconds())
}
func f6(outerFuncName string, f func(string) time.Duration, loop int, try int, i int) {
	defer wg.Done()
	el[loop][try][0][i] = fmt.Sprintf("%.02f", f(outerFuncName).Seconds())
}

func run(loop int, try int, gomaxprocs int) error {
	pprofFileName := fmt.Sprintf("binstat_external_test.loop-%d.try-%d.gomaxprocs-%d.pprof", loop, try, gomaxprocs)
	timerFile, err := os.Create(pprofFileName)
	if err != nil {
		return err
	}

	if err := pprof.StartCPUProfile(timerFile); err != nil {
		return err
	}

	// this function is purely for chewing CPU
	f := func(outerFuncName string) time.Duration {
		p := binstat.NewTime(outerFuncName, "")
		var sum int
		for i := 0; i < 10000000; i++ {
			sum -= i / 2
			sum *= i
			sum /= i/3 + 1
			sum -= i / 4
		}
		binstat.Dbg(p, fmt.Sprintf("%s() = %d", outerFuncName, sum))
		return binstat.End(p)
	}

	runtime.GOMAXPROCS(gomaxprocs)
	wg.Add(6)
	go f1("f-via-f1", f, loop, try, 0)
	go f2("f-via-f2", f, loop, try, 1)
	go f3("f-via-f3", f, loop, try, 2)
	go f4("f-via-f4", f, loop, try, 3)
	go f5("f-via-f5", f, loop, try, 4)
	go f6("f-via-f6", f, loop, try, 5)

	wg.Wait()
	pprof.StopCPUProfile()
	timerFile.Close()

	// run pprof and capture its output
	/*
		e.g. $ go tool pprof -top -unit seconds binstat_external_test.loop-1.try-2.gomaxprocs-8.pprof 2>&1 | egrep '(binstat_test.f|cum)'
		e.g.      flat  flat%   sum%        cum   cum%
		e.g.         0     0%   100%      0.07s 19.44%  github.com/onflow/flow-go/binstat_test.f1
		e.g.         0     0%   100%      0.02s  5.56%  github.com/onflow/flow-go/binstat_test.f2
		e.g.         0     0%   100%      0.06s 16.67%  github.com/onflow/flow-go/binstat_test.f3
		e.g.         0     0%   100%      0.11s 30.56%  github.com/onflow/flow-go/binstat_test.f4
		e.g.         0     0%   100%      0.06s 16.67%  github.com/onflow/flow-go/binstat_test.f5
		e.g.         0     0%   100%      0.03s  8.33%  github.com/onflow/flow-go/binstat_test.f6

		$ # todo: consider workaround: have seen pprof fail on macOS extremely infrequently, e.g. below .f5 completely missing?! how?!
		$ go tool pprof -top -unit seconds binstat_external_test.loop-1.try-2.gomaxprocs-8.pprof
		Type: cpu
		Time: Jun 14, 2021 at 7:37pm (PDT)
		Duration: 200.55ms, Total samples = 0.36s (179.51%)
		Showing nodes accounting for 0.36s, 100% of 0.36s total
				flat  flat%   sum%        cum   cum%
				0.36s   100%   100%      0.36s   100%  github.com/onflow/flow-go/binstat_test.run.func1
					0     0%   100%      0.07s 19.44%  github.com/onflow/flow-go/binstat_test.f1
					0     0%   100%      0.09s 25.00%  github.com/onflow/flow-go/binstat_test.f2
					0     0%   100%      0.06s 16.67%  github.com/onflow/flow-go/binstat_test.f3
					0     0%   100%      0.08s 22.22%  github.com/onflow/flow-go/binstat_test.f4
					0     0%   100%      0.06s 16.67%  github.com/onflow/flow-go/binstat_test.f6
	*/
	command := fmt.Sprintf("go tool pprof -top -unit seconds %s 2>&1 | egrep '(binstat_test.f|cum)'", pprofFileName)
	out, err := exec.Command("bash", "-c", command).Output()
	if err != nil {
		log.Fatal(err)
	}
	//debug fmt.Printf("debug: output of command: %s\n%s\n", command, out)

	// regex out the (cum)ulative column in pprof output
	r, _ := regexp.Compile(` ([0-9.]+)s`)
	matches := r.FindAllStringSubmatch(string(out), -1)
	//debug fmt.Printf("debug: matches=%#v\n", matches) // e.g. debug: matches=[][]string{[]string{" 0.04s", "0.04"}, []string{" 0.06s", "0.06"}, []string{" 0.08s", "0.08"}, []string{" 0.04s", "0.04"}, []string{" 0.09s", "0.09"}, []string{" 0.05s", "0.05"}}
	if len(matches) != funcs {
		panic(fmt.Sprintf("ERROR: INTERNAL: len(matches)=%d but should be %d; out=\n%s\n", len(matches), funcs, out))
	}

	// add the regex matches to a table of elapsed times
	for i := 0; i < funcs; i++ {
		//debug fmt.Printf("debug: matches[%d][1]=%s\n", i, matches[i][1])
		el[loop][try][1][i] = matches[i][1]
	}

	return nil
}

func TestWithPprof(t *testing.T) {
	os.Setenv("BINSTAT_ENABLE", "1")
	os.Setenv("BINSTAT_VERBOSE", "1")

	// delete any files hanging around from previous test run
	{
		command := fmt.Sprintf("ls -al ./binstat.test.pid-*.txt ./*gomaxprocs*.pprof ; rm -f ./binstat.test.pid-*.txt ./*gomaxprocs*.pprof")
		out, err := exec.Command("bash", "-c", command).Output()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("- debug: output of command: %s\n%s\n", command, out)
	}

	// run the test; loops of several tries running groups of go-routines
	for loop := 0; loop < loops; loop++ {
		gomaxprocs := 8
		if 0 == loop {
			gomaxprocs = 1
		}
		p := binstat.NewTime(fmt.Sprintf("loop-%d", loop), "")
		for try := 0; try < tries; try++ {
			fmt.Printf("- debug: loop=%d try=%d; running 6 identical functions with gomaxprocs=%d\n", loop, try+1, gomaxprocs)
			if err := run(loop, try, gomaxprocs); err != nil {
				log.Fatal(err)
			}
		}
		binstat.End(p)
	}

	// output a table of results similar to this
	/*
		- binstat------- pprof---------
		- try1 try2 try3 try1 try2 try3
		- 0.29 0.30 0.29 0.05 0.03 0.05 // f1() seconds; loop=0 gomaxprocs=1
		- 0.35 0.30 0.35 0.07 0.06 0.06 // f2() seconds; loop=0 gomaxprocs=1
		- 0.28 0.33 0.28 0.06 0.06 0.06 // f3() seconds; loop=0 gomaxprocs=1
		- 0.31 0.28 0.31 0.05 0.06 0.06 // f4() seconds; loop=0 gomaxprocs=1
		- 0.27 0.28 0.27 0.05 0.05 0.05 // f5() seconds; loop=0 gomaxprocs=1
		- 0.38 0.38 0.39 0.06 0.05 0.06 // f6() seconds; loop=0 gomaxprocs=1
		- binstat------- pprof---------
		- try1 try2 try3 try1 try2 try3
		- 0.07 0.07 0.07 0.05 0.03 0.07 // f1() seconds; loop=1 gomaxprocs=8
		- 0.07 0.07 0.07 0.05 0.04 0.03 // f2() seconds; loop=1 gomaxprocs=8
		- 0.07 0.07 0.07 0.04 0.07 0.07 // f3() seconds; loop=1 gomaxprocs=8
		- 0.07 0.07 0.07 0.05 0.02 0.08 // f4() seconds; loop=1 gomaxprocs=8
		- 0.07 0.07 0.07 0.09 0.06 0.07 // f5() seconds; loop=1 gomaxprocs=8
		- 0.07 0.07 0.07 0.04 0.10 0.03 // f6() seconds; loop=1 gomaxprocs=8
	*/
	for loop := 0; loop < loops; loop++ {
		fmt.Printf("- binstat------- pprof---------\n")
		fmt.Printf("-")
		for r := 0; r < 2; r++ {
			for try := 0; try < tries; try++ {
				fmt.Printf(" try%d", try+1)
			}
		}
		fmt.Printf("\n")
		gomaxprocs := 8
		if 0 == loop {
			gomaxprocs = 1
		}
		for i := 0; i < funcs; i++ {
			fmt.Printf("-")
			for mech := 0; mech < mechs; mech++ {
				for try := 0; try < tries; try++ {
					fmt.Printf(" %s", el[loop][try][mech][i])
				}
			}
			fmt.Printf(" // f%d() seconds; loop=%d gomaxprocs=%d\n", i+1, loop, gomaxprocs)
		}
	}

	// tell binstat to close down and write its stats file one last time
	binstat.Fin()

	// cat and sort binstat stats file
	{
		command := fmt.Sprintf("ls -al ./binstat.test.pid-*.txt ; cat ./binstat.test.pid-*.txt | sort")
		out, err := exec.Command("bash", "-c", command).Output()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("- debug: output of command: %s\n%s\n", command, out)
	}

	// todo: add more tests? which tests?

	// if we get here then no panic() or log.Fatal() calls happened :-)
	expected := 1
	actual := 1
	require.Equal(t, expected, actual)
}
