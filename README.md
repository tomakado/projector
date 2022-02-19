<p align="center">
	<img src="doc/logo.png" />
</p>

# Projector

A flexible, language and framework agnostic tool that allows you to generate projects from templates. 
Projector has some builtin templates, but you also can use your custom templates.

# Features

* Single binary, no extra dependencies
* Builtin templates that allow you to start quickly
* Simple template manifest format

# Installation

There are two ways to get Projector right now:

1. Get binary for your platform on [Releases page](https://github.com/tomakado/projector/releases)
2. Build Projector from source:
`go install github.com/tomakado/projector`

# Usage

## Projector CLI
Get general usage help with `-h` or `--help` flags:
```
❯ projector -h         
A flexible, language and framework agnostic tool that allows you to generate projects from templates. 
Projector has some builtin templates, but you can use your custom templates or import third-party templates
from GitHub.

Usage:
  projector [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  create      Create project using specified template
  help        Help about any command
  info        Show meta information about template
  list        List builtin and cached templates
  validate    Validate manifest without performing actions (dry run)

Flags:
  -h, --help      help for projector
  -v, --verbose   turn verbose mode on

Use "projector [command] --help" for more information about a command.
```

Create project with `create` command:
```
projector create go/hello-world --author "tomakado <hi@ildarkarymov.ru>" && \
								--name "my-awesome-app" && \
								--package "github.com/tomakado/go-helloworld" && \
							    ./hello-world/
```
where `go/hello-world` is template you want to use. You also can use custom manifest file by passing it with `--manifest` or `-m` flags.

List all available locally templates with `projector list`:
```
❯ projector list
go/hello-world
go/http
```

Get a bit more info about concrete template with `projector info [template]`:
```
❯ projector info go/hello-world
go/hello-world@1.0.0 by tomakado
URL: https://github.com/tomakado/projector
Description: Basic program to get started with Go
```

Validate custom manifest file with `projector validate --manifest=[path-to-custom-manifest-file]`:
```
❯ projector validate -m projector.toml
Manifest is valid ✅
```

Debug your template with `--verbose` flag:
```
❯ projector create go/hello-world -a "tomakado <hi@ildarkarymov.ru>" -n "my-awesome-app" -p "github.com/tomakado/go-helloworld" ./hello-world --verbose
2022/02/19 18:22:32 verbose mode is turned on
2022/02/19 18:22:32 using manifest name "go/hello-world" in embed fs
2022/02/19 18:22:32 initialized embedded fs provider
2022/02/19 18:22:32 working directory = "./hello-world"
2022/02/19 18:22:32 loading manifest "go/hello-world/projector.toml"
2022/02/19 18:22:32 [EmbedFSProvider] reading "go/hello-world/projector.toml" in "resources/templates/"
2022/02/19 18:22:32 parsing manifest
2022/02/19 18:22:32 validating manifest
2022/02/19 18:22:32 passing config and provider to new instance of *projector.Generator
2022/02/19 18:22:32 provider = *manifest.EmbedFSProvider, config = {WorkingDirectory:./hello-world ProjectAuthor:tomakado <hi@ildarkarymov.ru> ProjectName:my-awesome-app ProjectPackage:github.com/tomakado/go-helloworld Manifest:0xc0002a5960 ManifestPath:go/hello-world}
2022/02/19 18:22:32 initializing working directory "./hello-world"
2022/02/19 18:22:32 cd ./hello-world
2022/02/19 18:22:32 traversing manifest steps
2022/02/19 18:22:32 step "init go module and git repository", 1 of 2
2022/02/19 18:22:32 processing files
2022/02/19 18:22:32 extracting file template from "gitignore"
2022/02/19 18:22:32 [EmbedFSProvider] reading "go/hello-world/gitignore" in "resources/templates/"
2022/02/19 18:22:32 parsing file template
2022/02/19 18:22:32 rendering file
2022/02/19 18:22:32 saving rendered file to ".gitignore"
2022/02/19 18:22:32 parsing output path template ".gitignore"
2022/02/19 18:22:32 rendering output path template ".gitignore"
2022/02/19 18:22:32 mkdir .
2022/02/19 18:22:32 writing rendered file to ".gitignore"
2022/02/19 18:22:32 parsing shell script template "go mod init {{ .ProjectPackage }} && git init"
2022/02/19 18:22:32 rendering shell script
2022/02/19 18:22:32 executing shell script "go mod init github.com/tomakado/go-helloworld && git init"
2022/02/19 18:22:32 step "create project bootstrap", 2 of 2
2022/02/19 18:22:32 processing files
2022/02/19 18:22:32 extracting file template from "main.go.tpl"
2022/02/19 18:22:32 [EmbedFSProvider] reading "go/hello-world/main.go.tpl" in "resources/templates/"
2022/02/19 18:22:32 parsing file template
2022/02/19 18:22:32 rendering file
2022/02/19 18:22:32 saving rendered file to "main.go"
2022/02/19 18:22:32 parsing output path template "main.go"
2022/02/19 18:22:32 rendering output path template "main.go"
2022/02/19 18:22:32 mkdir .
2022/02/19 18:22:32 writing rendered file to "main.go"
```

## Manifest

Manifest is a file that describes how Projector should act to generate project. Projector uses [TOML](https://toml.io) format to store manifest on disk.

Example:
```toml
name="go/hello-world"
author="tomakado"
version="1.0.0"
url="https://github.com/tomakado/projector"
description="Basic program to get started with Go"

[[steps]]
name="init go module and git repository"
shell="go mod init {{ .ProjectPackage }} && git init"
	[[steps.files]]
	path="gitignore"
	output=".gitignore"
	
[[steps]]
name="create project bootstrap"
	[[steps.files]]
	path="main.go.tpl"
	output="main.go"
```

The first, top-level section containing fields `name`, `author`, `version`, `url`, `description` is about meta information. Then manifest must contain at least one _step._ Step must have name and either list of files to generate or shell script to execute.

### Reference

#### Manifest
| Field     | Description                                                                        |
| --------- | ---------------------------------------------------------------------------------- |
| `name`    | Name of template. Required.                                                        |
| `author`  | Author of template. Required.                                                      |
| `version` | Version of template. Required.                                                     |
| `url`     | URL of repository or website of template. Optional.                                |
| `steps`   | Array of steps. See [`step`](#step) for more info. Required at least one step. |

#### `step`
_Step_ is self-sufficient action performed by Projector to generate project. Projector “executes” steps sequentially, one by one. Inside `shell` field `text/template` syntax is supported, so you can use values exposed to [Template Context](#template-context) inside shell script.

| Field   | Description                                                                                                                                   |
| ------- | --------------------------------------------------------------------------------------------------------------------------------------------- |
| `name`  | Name of step. Required.                                                                                                                       |
| `files` | Array of files to generate. See [`file`](#file) for more info. Required if `shell` is not set.                                                |
| `shell` | Shell script to execute. `text/template` supported (see [Template Context](#template-context) for more info). Required if `files` is not set. |

#### `file`
_File_ in terms of Projector manifest is something like task of following kind:

1. Take file located at `path`;
2. Render file content as template with [Template Context](#template-context);
3. Render output path using `output` value;
4. Put rendered file to rendered `output` path.  

| Field    | Description                                                                                                                    |
| -------- | ------------------------------------------------------------------------------------------------------------------------------ |
| `path`   | Path to source file. `text/template` supported in content (see [Template Context](#template-context) for more info). Required. |
| `output` | Template of output path for rendered file. See [Template Context](#template-context) for more info. Required.                  |

#### Template Context
| Field              | Description                                                                             |
| ------------------ | --------------------------------------------------------------------------------------- |
| `WorkingDirectory` | Current working directory. Usually folder where project is being generated.             |
| `ProjectAuthor`    | Author of project.                                                                      |
| `ProjectName`      | Name of project.                                                                        |
| `ProjectPackage`   | Package name for project. E.g. in Go it would something like `github.com/owner/module`. |
| `Manifest`         | Reference to manifest. See [Manifest](#manifest) for info.                              |
# Backlog

There are some features I'd like to implement in Projector:

- [ ] Template for templates!
- [ ] Nice animated output of current step
- [ ] User-friendly error messages
- [ ] Import third-party templates from GitHub
- [ ] Custom options in template context
- [ ] Import third-party templates from any public or private git repository
- [ ] Template caching
- [ ] List local third-party cached templates
- [ ] Support for file masks on input files declaration
- [ ] `version` command (Go 1.18+)

# Contribution

If you want to help to develop Projector, the're several options:

* [Create issue](https://github.com/tomakado/projector/issues/new/choose) with problem you've met or any other kind
* Submit PRs with bug fixes or new features, especially from backlog :wink:
