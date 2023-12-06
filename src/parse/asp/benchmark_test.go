package asp

import (
	"arena"
	"bytes"
	"fmt"
	"github.com/thought-machine/please/rules"
	"github.com/thought-machine/please/src/core"
	"github.com/thought-machine/please/src/parse/asp/heap"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"
)

var code = `

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
	arena.NewArena().Free()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		parseInParallel(10, 5, false)
	}
}

func BenchmarkParseWithArena(b *testing.B) {
	arena.NewArena().Free()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		parseInParallel(10, 5, true)
	}
}

func parseInParallel(threads, repeats int, useArena bool) {
	wg := new(sync.WaitGroup)
	wg.Add(threads)
	pool := heap.NewPool(threads, -1, time.Second)
	for j := 0; j < threads; j++ {
		go func() {
			for k := 0; k < repeats; k++ {
				r := bytes.NewReader([]byte(code))
				if !useArena {
					parseFileInput(r, nil)
					continue
				}
				heap := pool.Get()
				parseFileInput(r, heap.Arena)
				pool.Return(heap)
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

func BenchmarkParseAndInterpretWithArena(b *testing.B) {
	b.ReportAllocs()
	p := newParserWithGo()
	arena.NewArena().Free()
	b.ResetTimer()

	parseAndInterpretInParallel(10, b.N*1000, true, p)
}

func BenchmarkParseAndInterpretWithoutArena(b *testing.B) {
	b.ReportAllocs()
	p := newParserWithGo()
	arena.NewArena().Free()
	b.ResetTimer()

	parseAndInterpretInParallel(10, b.N*1000, false, p)
}

func newParserWithGo() *Parser {
	p := NewParser(core.NewDefaultBuildState())
	bs, err := os.ReadFile("go.build_defs")
	if err != nil {
		panic(err)
	}

	// Dummy out the plugin config so the build defs can be interpreted
	goPluginConf := pyDict{
		"PKG_INFO":         pyBool(true),
		"IMPORT_PATH":      None,
		"LEGACY_IMPORTS":   pyBool(false),
		"RACE":             pyBool(false),
		"DEFAULT_STATIC":   pyBool(false),
		"BUILDMODE":        None,
		"STDLIB":           None,
		"C_FLAGS":          None,
		"LD_FLAGS":         None,
		"CGO_ENABLED":      None,
		"PLEASE_GO_TOOL":   pyString("please_go"),
		"GO_TOOL":          pyString("go"),
		"CC_TOOL":          pyString("cc"),
		"TEST_ROOT_COMPAT": pyBool(false),
	}

	p.interpreter.scope.config.base.dict["GO"] = goPluginConf

	assets, err := rules.AllAssets()
	if err != nil {
		panic(err)
	}

	for _, asset := range assets {
		bs, err := rules.ReadAsset(asset)
		if err != nil {
			panic(err)
		}

		if err := p.LoadBuiltins(asset, bs); err != nil {
			panic(err)
		}
	}

	if err := p.LoadBuiltins("go.build_defs", bs); err != nil {
		panic(err)
	}
	return p
}

func parseAndInterpretInParallel(threads, repeats int, withArena bool, p *Parser) {
	wg := new(sync.WaitGroup)
	wg.Add(threads)
	pool := heap.NewPool(threads, -1, time.Second)
	for thread := 0; thread < threads; thread++ {
		go func(thread int) {
			for repeat := 0; repeat < repeats; repeat++ {
				pkg := fmt.Sprintf("src/asp/parse_%v_%v", thread, repeat)
				var heap *heap.Heap
				var arena *arena.Arena
				if withArena {
					heap = pool.Get()
					arena = heap.Arena
				}

				s := p.interpreter.scope.newScope(core.NewPackage(pkg), arena, core.ParseModeNormal, filepath.Join(pkg, "BUILD"), 10)

				// This call to ParseData currently doesn't use the arena
				stmts, err := s.interpreter.parser.ParseData(arena, []byte(code), "BUILD")
				if err != nil {
					panic(err)
				}

				s.interpretStatements(stmts)

				// A quick assert to see if we have the target we expect
				t := s.state.Graph.TargetOrDie(core.NewBuildLabel(pkg, "asp"))
				_ = t
				if withArena {
					pool.Return(heap)
				}
			}
			wg.Done()
		}(thread)
	}

	wg.Wait()
}
