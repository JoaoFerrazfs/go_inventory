package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
)

const (
	ColorReset  = "\x1b[0m"
	ColorRed    = "\x1b[31m"
	ColorGreen  = "\x1b[32m"
	ColorYellow = "\x1b[33m"
)

type Event struct {
	Time    string `json:"Time"`
	Action  string `json:"Action"`
	Package string `json:"Package"`
	Test    string `json:"Test"`
	Output  string `json:"Output"`
}

func printColored(action, pkg, test string) {
	var color string
	var label string
	switch action {
	case "run":
		color = ColorYellow
		label = "[RUN ]"
	case "pass":
		color = ColorGreen
		label = "[PASS]"
	case "fail":
		color = ColorRed
		label = "[FAIL]"
	default:
		color = ColorReset
		label = "[INFO]"
	}
	fmt.Printf("%s%s %s %s%s\n\n", color, label, pkg, test, ColorReset)
}

func main() {
	// command-line flags
	summaryOnly := flag.Bool("summary-only", false, "print only the failing tests summary (no per-test lines)")
	flag.BoolVar(summaryOnly, "s", false, "shorthand for --summary-only")
	help := flag.Bool("help", false, "show help")
	flag.BoolVar(help, "h", false, "shorthand for --help")

	// custom usage/help message
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: <go test invocation> -json | pretty_print_tests [options]\n\n")
		fmt.Fprintf(os.Stderr, "Pretty prints `go test -json` output grouped by package.\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		fmt.Fprintf(os.Stderr, "  -s, --summary-only    Print only the failing tests summary (no per-test lines)\n")
		fmt.Fprintf(os.Stderr, "  -h, --help            Show this help message\n\n")
		fmt.Fprintf(os.Stderr, "Examples:\n")
		fmt.Fprintf(os.Stderr, "  go test ./... -json | /tmp/pretty_print_tests\n")
		fmt.Fprintf(os.Stderr, "  go test ./pkg -json | /tmp/pretty_print_tests --summary-only\n")
	}

	flag.Parse()
	if *help {
		flag.Usage()
		os.Exit(0)
	}

	scanner := bufio.NewScanner(os.Stdin)
	outputs := make(map[[2]string][]string)
	failures := make([][2]string, 0)
	pkgBuffers := make(map[string][][2]string) // pkg -> list of (action,test)
	pkgOrder := make([]string, 0)

	if !*summaryOnly {
		fmt.Println("Running tests (grouped by package):")
	}

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		var e Event
		if err := json.Unmarshal([]byte(line), &e); err != nil {
			continue
		}

		pkg := e.Package
		if pkg == "" {
			pkg = "<unknown>"
		}
		test := e.Test
		if test == "" {
			test = ""
		}

		key := [2]string{pkg, test}
		outputs[key] = append(outputs[key], e.Output)

		if _, ok := pkgBuffers[pkg]; !ok {
			pkgBuffers[pkg] = make([][2]string, 0)
			pkgOrder = append(pkgOrder, pkg)
		}

		if (e.Action == "run" || e.Action == "pass" || e.Action == "fail") && test != "" {
			// record events always; printing may be suppressed by summaryOnly
			pkgBuffers[pkg] = append(pkgBuffers[pkg], [2]string{e.Action, test})
			if e.Action == "fail" {
				// record failure once
				found := false
				for _, f := range failures {
					if f[0] == pkg && f[1] == test {
						found = true
						break
					}
				}
				if !found {
					failures = append(failures, key)
				}
			}
		}

		// package-level completion: flush the package buffer
		if test == "" && (e.Action == "pass" || e.Action == "fail") {
			// print buffered events for this package in insert order unless summaryOnly
			if !*summaryOnly {
				buf := pkgBuffers[pkg]
				for _, p := range buf {
					printColored(p[0], pkg, p[1])
				}
				// print separator
				fmt.Println(strings.Repeat("-", 120))
				fmt.Println()
			}
			delete(pkgBuffers, pkg)
		}
	}

	if len(failures) == 0 {
		if !*summaryOnly {
			fmt.Println("All tests passed (no failing tests detected).")
		} else {
			fmt.Println("All tests passed (no failing tests detected).")
		}
		os.Exit(0)
	}

	fmt.Println("\n==== Failing tests summary ====\n")
	for _, f := range failures {
		pkg := f[0]
		test := f[1]
		fmt.Printf("Package: %s  Test: %s\n\n", pkg, test)
		key := [2]string{pkg, test}
		for _, out := range outputs[key] {
			fmt.Print(out)
		}
		fmt.Println("\n----\n")
	}

	os.Exit(1)
}
