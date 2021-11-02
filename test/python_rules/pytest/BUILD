package(python_test_runner = "pytest")

python_library(
    name = "inc",
    srcs = ["inc.py"],
)

python_library(
    name = "relative_imports",
    srcs = glob(["relative_imports/*"]),
)

python_test(
    name = "test_relative_imports",
    srcs = ["test_relative_imports.py"],
    deps = [
        ":relative_imports",
    ],
)

python_test(
    name = "pytest_test",
    srcs = ["test_runner.py"],
    deps = [":inc"],
)

python_test(
    name = "test_pytest_unicode",
    srcs = ["test_pytest_unicode.py"],
    deps = [
        ":pytest_file",
        "//third_party/python:pytest",
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

python_test(
    name = "testytest",
    srcs = ["quicktest.py"],
)

filegroup(
    name = "pytest_file",
    srcs = ["pytest.py"],
)
