summary: Getting started with plugins
description: Getting a plugin set up for your Please repository
id: plugin_setup
categories: beginner
tags: medium
status: Published
authors: Sam Westmoreland
Feedback Link: https://github.com/thought-machine/please

# Getting started with Plugins
## Overview
Duration: 4

### Prerequisites
- You must have Please installed: [Install please](https://please.build/quickstart.html)

### What You’ll Learn
- How to write your own custom plugin for use with Please

### What if I get stuck?

The final result of running through this codelab can be found
[here](https://github.com/thought-machine/please-codelabs/tree/main/getting_started_python) for reference. If you really
get stuck you can find us on [gitter](https://gitter.im/please-build/Lobby)!

## Initialising your project
Duration: 2

Let's create a new project:
```shell
$ mkdir my_plz_plugin && cd my_plz_plugin
$ plz init --no_prompt
```

## Setting up your config file
Duration: 3

We now have an initialised Please repo. Let's set the config file up so that `plz` will know that this is a plugin repo:

### `.plzconfig`
```ini
[PluginDefinition]
name = fooplug
```
Once we've done that, we can add our plugin-specific config values that will allow any users of our plugin to configure the plugin as they need.

### `.plzconfig`
```ini
[PluginConfig "hash_function"]
ConfigKey = HashFunction
DefaultValue = SHA256
Optional = false
```
The ConfigKey can be set or not. If we choose not to specify it, it will be inferred from the section title (e.g. `hash_function` implies `HashFunction`).

## Writing a new build definition
Duration: 3

Now let's add a simple build definition. For more information on writing custom build definitions, have a look at the ['Custom build rules with genrule()'](https://please.build/codelabs/genrule/) codelab.

### `./build_defs/foolang.build_defs`
```python
def hash_file(name:str, file:str) -> str:
    return genrule(
        name = name,
        srcs = [file],
        outs = [f'{name}.hash'],
        cmd = 'sha1sum $SRC > $OUT',
    )
```
We can then expose the build definition using a filegroup like so:
### `build_defs/BUILD`
```python
# Expose our foolang build_defs
filegroup(
    name = "foolang",
    srcs = ["foolang.build_defs"],
    visibility = ["PUBLIC"],
)
```

## Testing our plugin
To test our new build definition, we'll create a dummy file:
### `test/somefile.foo`
```
foo
bar
baz
```

### `test/BUILD`
```python
subinclude("//build_defs:foolang)

hash_file(
  name = "hashme",
  file = "somefile.foo",
)
```

And then we'll make sure it works:

```shell
$ plz build //test:hashme
Build finished; total time 80ms, incrementality 0.0%. Outputs:
//test:hashme:                                                
  plz-out/gen/test/hashme.hash                                

$ cat plz-out/gen/test/hashme.hash
0562f08aef399135936d6fb4eb0cc7bc1890d5b4  test/somefile.foo
```

## Testing our code
Duration: 5

Let's create a very simple test for our library:
### `src/greetings/greetings_test.py`
```python
import unittest
from src.greetings import greetings

class GreetingTest(unittest.TestCase):

    def test_greeting(self):
        self.assertTrue(greetings.greeting())

```

We then need to tell Please about our tests:
### `src/greetings/BUILD`
```python
python_library(
    name = "greetings",
    srcs = ["greetings.py"],
    visibility = ["//src/..."],
)

python_test(
    name = "greetings_test",
    srcs = ["greetings_test.py"],
    # Here we have used the shorthand `:greetings` label format. This format can be used to refer to a rule in the same
    # package and is shorthand for `//src/greetings:greetings`.
    deps = [":greetings"],
)
```

We've used `python_test()` to define our test target. This is a special build rule that is considered a test. These
rules can be executed as such:
```
$ plz test //src/...
//src/greetings:greetings_test 1 test run in 3ms; 1 passed
1 test target and 1 test run in 3ms; 1 passed. Total time 90ms.
```

Please will run all the tests it finds under `//src/...`, and aggregate the results up. This works even across
languages allowing you to test your whole project with a single command.

## Third-party dependencies
Duration: 7

### Using `pip_library()`

Eventually, most projects need to depend on third-party code. Let's include NumPy into our package. Conventionally,
third-party dependencies live under `//third_party/...` (although they don't have to), so let's create that package:

### `third_party/python/BUILD`
```python
package(default_visibility = ["PUBLIC"])

pip_library(
    name = "numpy",
    version = "1.18.4",
    zip_safe = False, # This is because NumPy has shared object files which can't be linked to them when zipped up
)
```

This will download NumPy for us to use in our project. We use the `package()` built-in function to set the default
visibility for this package. This can be very useful for third-party rules to avoid having to specify
`visibility = ["PUBLIC"]` on every `pip_library()` invocation.

NB: The visibility "PUBLIC" is a special case. Typically, items in the visibility list are labels. "PUBLIC" is equivalent
to `//...`.

### Setting up our module path
Importing Python modules is based on the import path. That means by default, we'd import NumPy as
`import third_party.python.numpy`. To fix this, we need to tell Please where our third-party module is. Add the
following to your `.plzconfig`:

### `.plzconfig`
```
[python]
moduledir = third_party.python
```

### Updating our tests

We can now use this library in our code:

### `src/greetings/greetings.py`
```go
from numpy import random

def greeting():
    return random.choice(["Hello", "Bonjour", "Marhabaan"])

```

And add NumPy as a dependency:
### `src/greetings/BUILD`
```python
python_library(
    name = "greetings",
    srcs = ["greetings.py"],
    visibility = ["//src/..."],
    deps = ["//third_party/python:numpy"],
)

python_test(
    name = "greetings_test",
    srcs = ["greetings_test.py"],
    deps = [":greetings"],
)
```

## What next?
Duration: 1

Hopefully you now have an idea as to how to build Python with Please. Please is capable of so much more though!

- [Please basics](/basics.html) - A more general introduction to Please. It covers a lot of what we have in this
tutorial in more detail.
- [Built-in rules](/lexicon.html#python) - See the rest of the Python rules as well as rules for other languages and tools.
- [Config](/config.html#python) - See the available config options for Please, especially those relating to Python.
- [Command line interface](/commands.html) - Please has a powerful command line interface. Interrogate the build graph,
determine files changes since master, watch rules and build them automatically as things change and much more! Use
`plz help`, and explore this rich set of commands!

Otherwise, why not try one of the other codelabs!

