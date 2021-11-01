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
