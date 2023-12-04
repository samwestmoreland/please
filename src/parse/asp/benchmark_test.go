package asp

import (
	"bytes"
	"sync"
	"testing"
	"arena"
)

var code = `subinclude("//build_defs:benchmark")

go_library(
    name = "asp",
    srcs = [
        "builtins.go",
        "config.go",
        "errors.go",
        "exec.go",
        "file_position.go",
        "grammar.go",
        "grammar_parse.go",
        "interpreter.go",
        "lexer.go",
        "objects.go",
        "parser.go",
        "targets.go",
        "util.go",
    ],
    pgo_file = "//:pgo",
    visibility = ["PUBLIC"],
    deps = [
        "///third_party/go/github.com_Masterminds_semver_v3//:v3",
        "///third_party/go/github.com_manifoldco_promptui//:promptui",
        "///third_party/go/github.com_please-build_gcfg//types",
        "///third_party/go/golang.org_x_exp//slices",
        "//rules",
        "//src/cli",
        "//src/cli/logging",
        "//src/cmap",
        "//src/core",
        "//src/fs",
    ],
)

go_test(
    name = "asp_test",
    srcs = [
        "builtins_test.go",
        "config_test.go",
        "file_position_test.go",
        "interpreter_test.go",
        "label_context_test.go",
        "lexer_test.go",
        "logging_test.go",
        "parser_test.go",
        "targets_test.go",
        "util_test.go",
    ],
    data = ["test_data"],
    deps = [
        ":asp",
        "///third_party/go/github.com_stretchr_testify//assert",
        "///third_party/go/github.com_stretchr_testify//require",
        "///third_party/go/gopkg.in_op_go-logging.v1//:go-logging.v1",
        "//rules",
        "//src/cli",
        "//src/core",
    ],
)
`

func BenchmarkParse(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		parseInParallel(10, 5, false)
	}
}

func BenchmarkParseWithArena(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		parseInParallel(10, 5, true)
	}
}

func parseInParallel(threads, repeats int, useArena bool) {
	wg := new(sync.WaitGroup)
	wg.Add(threads*repeats)
	for j := 0; j < threads; j++ {
		go func() {
			for k :=0 ; k < repeats; k++ {
				r := bytes.NewReader([]byte(code))
				if !useArena {
					parseFileInput(r, nil)
					wg.Done()
					continue
				}
				a := arena.NewArena()
				parseFileInput(r, a)
				a.Free()
				wg.Done()
			}
		}()
	}
	wg.Wait()
}