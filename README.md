# go-import-name-validator
Enforces user defined semantics around import naming.

## Usage

```
# Install from local repo
go install $GOPATH/src/github.com/sublime-security/go-import-name-validator/cmd/inspect-imports

# Install from GitHub
go install github.com/sublime-security/go-import-name-validator/cmd/inspect-imports@main
```

`inspect-imports` should be available in `$GOPATH/bin`, which is assumed to be on your `PATH`.

Your working director must be the repo you're inspecting. Example to inspect a package within this repo (w/o any requirements specified):

```
cd $GOPATH/src/github.com/sublime-security/go-import-name-validator
inspect-imports $GOPATH/src/github.com/sublime-security/go-import-name-validator/imports_analyzer
```

Similarly to check all packages in the repo

```
cd $GOPATH/src/github.com/sublime-security/go-import-name-validator
go list ./... | xargs inspect-imports
```

All standard Go singlechecker flags are supported. Such as `-fix` to fix inline, where possible.
```
inspect-imports -fix $GOPATH/src/github.com/sublime-security/go-import-name-validator/imports_analyzer
```

### Required Package Names

The `-require-name` flag may be specified any number of times to require that:
* The given name is only used for the given path
* Whenever the given path is imported, the name is used.

`-require-name <path>=<?name>`

Examples:
* `-require-name "github.com/sirupsen/logrus"=log`
  * The logrus package must be imported as `log`
  * If the name `log` is used, it must be for `"github.com/sirupsen/logrus"`
* `-require-name "errors"=`
  * If the errors package is imported, it must be unnamed (it cannot be named `errors`)

### Forbidden Imports

The `-forbidden` flag prohibits an import path from being included, any may be specified any number of times.

Examples:
* `-forbidden github.com/RichardKnop/machinery/v2/log`
  * `github.com/RichardKnop/machinery/v2/log` may not be imported
  * No partial matching is supported, although this could be a good addition.

### Full Examples

```
cd $GOPATH/src/github.com/sublime-security/go-import-name-validator
go list ./... | xargs inspect-imports -fix -forbidden github.com/RichardKnop/machinery/v2/log -require-name "github.com/sirupsen/logrus"=log -require-name "errors"=

# Same but single package
inspect-imports -fix -forbidden github.com/RichardKnop/machinery/v2/log -require-name "github.com/sirupsen/logrus"=log -require-name "errors"= $GOPATH/src/github.com/sublime-security/go-import-name-validator/imports_analyzer
```