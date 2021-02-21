# Jen

> Jen - _noun_ (in Chinese philosophy) a compassionate love for humanity or for the world as a whole.

Jen is a project scaffolding and script runner that accompanies your project throughout its lifetime.

# Motivation

We were not satisfied with existing project scaffolding tools, which often are language-specific and
leave you on your own once your project has been generated. Many DevOps shell scripts (ie:
registering your project with CI/CD and other infra, promoting it from staging to prod...) need to be
maintained, organized and shared separately. They typically require similar inputs to those provided
during scaffolding (ie: project name, cluster, cloud region, team...), yet you have to pass that
information as arguments all over again every time you invoke them.

As DevOps, we have many concerns to address, such as:

- How do we organize and share reusable project templates?
- What about project-specific and common DevOps shell scripts?
- How do we deal with those scripts' arguments/variables?

Jen aims to provide a very simple framework to answer all those questions.

# How it works

- Put all you project templates in a git repo with a specific structure, including shell scripts
  specific to each template or common to all.
- Jen automatically clones that repo locally upon first use.
- When scaffolding a project, Jen prompts user for which template to use and all values required by
  template and associated shell scripts
- All that information gets stored in a `jen.yaml` file in the project's root dir and
  remains available everytime you run a `jen ...` command within your project's directory structure.
- Use jen throughout your project's life to run companion shell scripts using same variables.

# Getting started

## Install Jen

## Create git repo for templates/scripts

Create a git repo to store all your templates and scripts, using the following structure:

## Set Jen env vars

- JEN_HOME: Directory where Jen will clone your templates git repo (defaults to `~/.jen`)
- JEN_REPO: URL of your templates git repo to clone

## Creating a new project from a template

```bash
$ mkdir foobar
$ cd foobar
$ jen do create
```

1. If it doesn't already exist, Jen will automatically clone your templates git repo into your $JEN_HOME.
2. It will prompt you to select from the list of templates available in that repo. That template name
   will be stored in `jen.yaml` file in current dir.
3. Depending on the steps configured in the `create` action of the template you selected, you should
   typically be prompted for the values of different variables, which will also be saved in the `jen.yaml`
   file.
4. Typically, the selected template will be rendered to current directory and, potentially, some scripts
   may automatically run to register your project with your infrastructure.

# Jen commands

## Invoke an action

To invoke given action from your project's template spec:

```bash
$ jen do ACTION
```

All steps defined within that action will be called in order.

## Execute a shell command

To execute any shell command, including your custom shell scripts, with your project's variables as environment:

```bash
$ jen exec COMMAND ARG1 ARG2 ...
```

## Start a sub-shell

To start a sub-shell with your custom shell scripts added to `$PATH` and your project's variables as environment:

```bash
$ jen exec SHELL
```

Where `SHELL` can any of `bash`, `zsh`, `sh`...

You are then free to call as many custom scripts and shell commands as you want, until you do `exit`.

# Update templates repo

To pull latest version of templates git repo:

```bash
$ jen pull
```

# Templates git repo

- bin (scripts common to all templates)
- templates
  - TEMPLATE_NAME
    - spec.yaml (defines actions/steps/variables)
    - src (template files to render)
    - bin (template-specific scripts)

# `spec.yaml` files

Each template has a `spec.yaml` file in its root that specifies how to render the template, what
variables to prompt user and what actions user can invoke throughout the project's lifetime.

It has this general structure:

```yaml
version: 0.2.0
description: ...
actions:
  ACTION1:
    - STEP1: ...
    - STEP2: ...
    - STEP3: ...
  ACTION2:
    - STEP1: ...
```

## Actions

Actions are named operations that can be invoked by user via the `jen do ACTION` command,
or as part of another action. They can have arbitrary names, however it is recommended to
follow the convention having at least the following two actions:

- `create`: action that initially scaffolds the project
- `prompt`: action that prompts user for variables (this action is typically invoked from
  the `create` action)

The order of actions is irrelevant, much like the definition of functions in any program.

## Steps

Each action is comprised of one or many steps that are executed sequentially when the
action is invoked (their order is therefore important).

Steps have predefined names and purposes:

- `if`: conditionally invokes child steps
- `do`: executes another action by name (much like a function call)
- `exec`: executes a shell command, including custom scripts, with project vars in environment
- `render`: renders template into current dir, using project vars
- `input`: prompts user for a single free-form string var
- `choice`: prompts user for a single string var among a list of multiple proposed choices
- `option`: prompts user for single boolean var as a yes/no question
- `options`: prompts user for multiple boolean vars as a list of toggles

## Example

```yaml
# Version of jen file format (for future compatibility checks)
version: 0.2.0

# Description displayed to user during template selection
description: The customary Hello World example

# Actions are sets of steps that can be invoked by user by their name
actions:
  # By convention, the "create" action is in charge of scaffolding the project initially
  create:
    # This step invokes the "prompt" action defined below
    - do: prompt
    # This step renders the "./src" template sub-dir into current dir
    - render: ./src
    # This step only executes its child steps if given expression evaluates to true
    - if: .INSTALL
      then:
        # This step invokes the "install" action defined below
        - do: install

  # By convention, the "prompt" action is in charge of prompting user for project
  # variables. It is typically invoked as the first step of "create" action above, but
  # can also be invoked manually by user at a later time to modify variables or to
  # associate a template with an existing project that was not initially generated by
  # jen.
  prompt:
    # The "input" step prompts user for a single string variable
    - input:
        question: Project name
        var: PROJECT
        # Here we use the special "projectDirName" variable to propose to the user the
        # project's directory name as default project name, which is most often the case.
        default: "{{ .projectDirName }}"

    # The "option" step prompts user for a single boolean variable as a yes/no question
    - option:
        question: Do you want to register project {{ .PROJECT }} in infrastructure?
        var: INSTALL
        default: true

    # The "options" step prompts user for multiple boolean variables as a list of toggles
    - options:
        question: Select desired features
        items:
          - text: PostgreSQL database
            var: PSQL
            default: true
          - text: NewRelic instrumentation
            var: NEWRELIC
            default: false

    # The "choice" step prompts user for a single string value from a list of choices.
    - choice:
        question: What is your team?
        var: TEAM
        default: backend
        items:
          - text: Back End
            value: backend
          - text: Front End
            value: frontend
          - text: DevOps
            value: devops
          - text: Site Reliability Engineering
            value: sre
          - text: Customer Success Engineering
            value: cse

  # By convention, the "install" action is in charge of setting the project up with CI/CD
  # and infra. It is typically invoked as the last step of the "create" action.
  install:
    # The "exec" step allows to invoke shell commands, including custom scripts, while
    # passing them all project variables as env vars. In this case, it specifies a list
    # of multiple commands to execute.
    - exec:
        - create-container-repo
        - create-cicd-triggers

  # By convention, the "uninstall" action is in charge of removing the project from infra
  uninstall:
    # Here the "exec" step is invoked multiple times, each executing a single command
    - exec: remove-container-repo
    - exec: remove-cicd-triggers
```

# Templates

## Go template language

Jen leverages the Go templating engine described [here](https://golang.org/pkg/text/template/). It
also augments Go's built-in functions with the [sprig](https://masterminds.github.io/sprig/) set of
very helpful functions.

Those template expressions can be used in templates, user prompts and file and directory names, as
described in the following sections.

## Activating/deactivating rendering

By default all files in a template are copied as is, without rendering their content as templates.
Template rendering can however be activated or deactivate selectively on a per-file or per-directory
basis, by appending a `.tmpl` or `.notmpl` extension to file and directory names. Applying those
extensions to a directory affects all child files recursively, unless overriden down the tree.

Note that the `.tmpl` and `.notmpl` extensions are automatically stripped away from target file and
directory names.

To override the no templating default, you can simply append a `.tmpl` extension to the directory
name passed to the `render` step:

```yaml
- render: ./src.tmpl
```

## Dynamic file and directory names

File and directory names can include template expressions:

- src
  - {{.PROJECT}}.txt
  - {{.TEAM}} team files
    - file1.txt
    - file2.txt

## Conditional files and directories

Files and directories can be selectively included/excluded by including a double-square-bracket expression
in their name, which must evaluate to true in order for the file or directory to be included in render.

For example, in the following template directory structure:

- src
  - database[[.DB]]
    - migration.go
    - driver.go

the `database[[.DB]]` directory and its content will only be rendered if the `DB` var is `true` and the
double-square-bracket expression will automatically get stripped away from the target name:

- src
  - database
    - migration.go
    - driver.go

## Expressions in prompts

In prompt steps (`input`, `choice`, `option`, `options`) you can use template expressions within messages,
proposed choices and default values, by enclosing those expressions between `{{` and `}}`.

## Expressions in `if` step

The conditional for `if` steps is always a template expression, so _do not_ enclose it between double
braces:

```yaml
- if: .INSTALL
  then:
    - ...
```

## Placeholders

Because the PROJECT variable is typically used pervasively throughout templates in the form of `{{.PROJECT}}`
and `{{.PROJECT | upper}}`, we have introduced the special placeholders `projekt` and `PROJEKT`, which can
be used anywhere in file/dir names and templates without double-braces.

For example, the text "MY*PROJEKT_FILE" is equivalent to "MY*{{.PROJECT | upper}}\_FILE".

Currently, those two placeholders are hardcoded and are the only ones supported, but we plan to add support
for defining your own in the template spec.

# Shell scripts

When executing any action or shell command, jen automatically adds all project variables to your
shell environment. It also adds the `bin` directories for your template-specific and common custom scripts
to your `$PATH` env var. That way, your custom scripts can only be invoked when all proper environment
variables are set.

## Wishlist

- Invoking `jen do` without specifying an action should prompt user to select it from available list of actions.
- Add `jen export` command to output env variables in a format that can be sourced directly.
- Add reusable modules (including both templates and scripts).
- Add `set` step to set multiple variables.
- Add `confirm` step (similar to `if`, but `confirm` property contains message to display and `then` the steps to execute)
- Add `--dry-run` flag (automatically turns on `--verbose`?)
- Add regex validation for `input` prompt.
- Allow `do` step to define multiple actions to call.
- Allow to customize placeholders in spec file:

```
placeholders:
  projekt: {{.PROJECT | lower}}
  PROJEKT: {{.PROJECT | upper}},
```
