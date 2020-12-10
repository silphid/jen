# Jen

> Jen - *noun* (in Chinese philosophy) a compassionate love for humanity or for the world as a whole.

Code generator for scaffolding microservices from templates boasting best practices.

## Problem

We have no automated way to create new micro-services and integrate them into our CI/CD pipeline.

It is overly long and complex to manually create a new micro-service and best practices are often not followed or forgotten.

## Objectives

- Reduce project skeleton creation time from 1 day to 5 minutes:
  - Engineers being able to create a brand new service with hello world functionality, boasting our most important best practices, fully integrated into our CI/CD infrastructure, built and deployed in under 5 minutes.
- Promote documenting and implementing industry/team best practices
  - Logging and observability
  - Configuration (12 factor app)
  - Security
  - ...
- Generates project based on templates for different languages and project archetypes
- Fully customizable script per template
  - Custom prompts
    - Strings
    - Options
    - Selection
  - Custom actions (automatically & manually executable)
  - Conditional steps and template files based on selected options
  - Reproducible recipe (save values)
  - Modularity/reusability of template parts

## Features

- Custom prompts at command line
- Save answers to yaml file
- Create project directory structure from go templates

## Inputs

- Template to use (language / variant)
- Project name
- Custom vars defined in template

## Resources

- Yaml parsers:
  - https://sweetohm.net/article/go-yaml-parsers.en.html
  - https://github.com/kylelemons/go-gypsy

- My script drafts:
  - [create-service](https://github.com/Samasource/factotum/blob/master/rootfs/root/bin/pipelines/create-service)
  - [add-service-triggers](https://github.com/Samasource/factotum/blob/master/rootfs/root/bin/pipelines/add-service-triggers)
  - [remove-service-triggers](https://github.com/Samasource/factotum/blob/master/rootfs/root/bin/pipelines/remove-service-triggers)

## How it works

### Creating new project

- Create an empty directory and chdir into it
- Call `jen create {template}` (or just `jen create` to be prompted for the template)
- Select which template to use
- Answer template-specific questions

### Executing custom actions

- chdir into project directory
- Call `jen do {action}` (or just `jen do` to be prompted for the action)

### Start a shell with all project env vars exported

- cd into project directory
- Call `jen shell`



# Command syntax

## Create project in current empty-folder

``` bash
$ mkdir my-project
$ cd my-project
$ jen create
```

- Prompts for which template to use
- Prompts for all custom inputs defined in template
- Saves values to `jen.yaml` file in project root
- Copies and interpolates all source files

## Use values from yaml

``` bash
$ jen create -f jen.yaml
```

## Specify value (avoid prompt)

``` bash
$ jen create --set template=go-service
```

## Perform an action (ie: add/remove Codefresh triggers)

``` bash
$ jen do install
```
Manually invokes an action defined in spec file, using project values stored in jen.yaml file

## Dry run (only show output, skip all disk changes)

``` bash
$ jen create --dry-run
```

## Edit existing values

``` bash
$ jen edit
```

# Yaml formats

## Template spec

``` yaml
name: go-service
description: Generic Go micro-service implementing all best practices
create:
- value:
    name: NAME
    title: Name of project # defaults to value of "name" property above
    default: $(basename $(pwd))
- option:
    name: myOption
    title: A checkbox choice
    steps: # optional sub-steps (if any sub-values define, value of myOption itself gets stored in myOption.enabled)
    - value: # only prompted if parent option selected
        name: subValue
        title: Sub value of myOptions
- multi:
    title: Select all desired options
    items: # Multi-selection (check-boxes), sets variable "myOption1" to either "true" or "false", and same thing for "myOption2"
    - name: myOption1
    	title: Displayed text for option 1
    - name: myOption2
    	title: Displayed text for option 2
- select:
    name: mySelect
    title: Select one of the following choices
    items: # Exclusive selection (radio-button), sets variable "Choices" to one of "choice1" or "choice2"
    - value: choice1
    	title: Displayed text for choice 1
    - value: choice2
    	title: Displayed text for choice 2
- option: # name is omitted here, because not needed
    title: Install Codefresh triggers?
    steps:
    - do: install # All execution is enqueued to be executed at the end
- render: src
- if: mySelect == choice1
  render: ../infra
- if: install
  do: install
- exec: echo "Hello world"
actions:
  install:
  - if: myOption == some value
  	exec: install # executes script "install", passing all values as env vars
  uninstall:
  - exec: uninstall # executes script "uninstall", passing all values as env vars
```

## Expressions

The `title`, `value`, `default` and `if` properties are treated as expressions, which are strings containing either or both go template and shell directives.

For example, `title` could be defined as:

```yaml
title: Do you want your project {{.NAME}} to be created in current directory $(basename $(pwd))
```

## Values file

``` yaml
template: go-service
name: MyProject
createGit: true
install: true
key1: value1
key2: value2
key3: value3
```

## Templates are simply go templates

```handlebars
Any text {{ .myValue }} can be templated!
```

- Support for the whole [sprig](https://github.com/Masterminds/sprig) library

## Dynamic file names

```handlebars
FileName{{.myValue}}.txt
```

## Conditional folders

```handlebars
FolderName[[.myValue]]
```

Folder will be rendered only if expression within double-brackets evaluates to true.

Double-bracket part can be specified anywhere within name (start, middle or end) and will anyway be stripped away from file name.

If only double-bracket part is specified (no folder name), the entire folder's content gets promoted one level up.

## Conditional files

```
FileName[[.myValue]].txt
```

File will be rendered only if expression within double-brackets evaluates to true.

Double-bracket part can be specified anywhere within name (start, middle or end) and will anyway be stripped away from file name.

# Go Best practices

## Libraries

- Environment vars: https://github.com/kelseyhightower/envconfig
- Logging: https://github.com/uber-go/zap

## Boilerplates

- podinfo template: https://github.com/stefanprodan/podinfo
- https://github.com/dhax/go-base

## Practices

- Unit tests
- Documentation
- OpenAPI
- [Organizing Go Code](https://blog.golang.org/organizing-go-code)
- [Standard Go Project Layout](https://github.com/golang-standards/project-layout)
- Codefresh build badge

# Pseudo-code / logic

## Template processing

- Fill values from provided values file (`--file jen.yaml`) into dict
- Fill command-line provided values (`--set myValue=value`) into dict
- Execute all steps in order
  - Skip steps for which value is already provided
  - Execution of commands is enqueued and delayed until the end of all interactive prompts

# Packaging

## jen

- Source code only

## factotum

- Built and packaged as part of `factotum` docker image
- All templates
- Custom action scripts
  - Create git repo
  - Create Codefresh triggers
  - Remove Codefresh triggers
  - ...

# Backlog

- ~~Create project skeleton~~
  - ~~Cobra~~
  - ~~Viper~~
- ~~Load spec file~~
  - ~~Model spec file (dict/structs)~~
  - ~~Load model from spec file~~
- ~~Prompts~~
- ~~Rendering~~
- ~~Add support for [sprig](https://github.com/Masterminds/sprig) functions~~
- ~~Dynamic folder/file names~~
- ~~Conditional folders/files~~
- ~~If (conditionals)~~
- ~~Do (actions)~~
- ~~Exec (shell)~~
  - ~~Pass all values as env vars (ie: `myValue` becomes `MY_VALUE`)~~
- ~~Select template from list~~
- ~~Configure template path in `~/.jen` config file~~
- ~~Automatically create output dir based on project name~~
- ~~Execute shell commands in output folder (instead of CWD)~~
- ~~Factotum~~
  - ~~Build and configure~~
  - ~~Create first draft of go template~~
  - ~~Create Codefresh build triggers~~
  - ~~Create git repo~~

## Wishlist (if time allows)

- Reusable modules
- Add `set` step to set multiple variables (are those saved to `jen.yaml`?)
- Override values at command-line level (`--set myValue=value`)
- Factotum Dockerfile should specify explicit jen version number in `go get`
- Escape templating:
    - Only treat as templates files ending with `.gotmpl` (and remove extension)
    - Exceptionally escape within any file, using `{{{` and `}}}` to represent `{{` and `}}`
- Per-template/module scripts in `bin` dir, which are automatically included in `PATH`
- `confirm` step (similar to `if`, but `confirm` property contains message to display and `then` the steps to execute)
- Custom placeholders:
```
placeholders:
  projekt: {{.PROJECT | lower}}
  PROJEKT: {{.PROJECT | upper}}
```

# Structure of JEN_HOME directory

- bin
- modules
  - NAME
    - spec.yaml
    - template
    - bin
- templates
  - NAME
    - spec.yaml
    - src
    - bin

# Innovation Week 2.0

- Finish refactoring/clean-up/unit tests
- Baked-in placeholders
- `JEN_HOME` env var to define templates location?  
- Upon `jen do <action>`, automatically find and load existing `jen.yaml` if any
  - File should contain:
    - Template name used
    - Input values
  - If file not found, prompt for template and then save `jen.yaml` with selected template.
- Each prompt should immediately save all values to `jen.yaml`
- `jen do <action>`
- Use `kaguya` service to create first template
- Automate `jen` builds/releases
- `shell`
- `exec`
- `export`


