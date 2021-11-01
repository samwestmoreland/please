build_defs = glob(["*.build_defs"])

go_library(
    name = "rules",
    srcs = ["rules.go"],
    resources = build_defs,
    visibility = ["PUBLIC"],
)

go_test(
    name = "rules_test",
    srcs = ["rules_test.go"],
    resources = build_defs,
    visibility = ["PUBLIC"],
    deps = [
        ":rules",
        "//third_party/go:testify",
    ],
)

filegroup(
    name = "pytest_deps",
    visibility = ["PUBLIC"],
    srcs = [
        "pytest.py",
        # "//third_party/python:pytest",
        "//third_party/python:py",
        "//third_party/python:pluggy",
        "//third_party/python:attrs",
        "//third_party/python:funcsigs",
        "//third_party/python:more_itertools",
        "//third_party/python:packaging",
        "//third_party/python:pkg_resources",
        "//third_party/python:importlib_metadata",
        "//third_party/python:zipp",
    ],
)

filegroup(
    name = "behave_deps",
    visibility = ["PUBLIC"],
    srcs = [
        "behave.py",
        "//third_party/python:parse",
        #"//third_party/parse.py", # this was the original
        "//third_party/python:behave",
        "//third_party/python:parse_type",
        "//third_party/python:traceback2",
        "//third_party/python:enum34",
        "//third_party/python:win_unicode_console",
        "//third_party/python:colorama",
    ],
)
