# Jen

> Jen - _noun_ (in Chinese philosophy) a compassionate love for humanity or for the world as a whole.

Jen is a project scaffolding and script runner that accompanies your project throughout its life-time.

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
  template and associated shell scripts.
- All that information gets stored in a `jen.yaml` file in the project's root dir and
  remains available everytime you run a `jen ...` command anywhere within your project's directory
  structure.
- Use jen throughout your project's life-time to run companion shell scripts using same variables.

# Getting started

## Install Jen

Download and install latest release from [github](https://github.com/Samasource/jen/releases/latest).

## Create git repo for templates/scripts

Create a git repo to store all your templates and scripts, using the following structure:

- `bin` (scripts common to all templates)
- `templates`
  - `TEMPLATE_NAME`
    - `spec.yaml` (defines actions/steps/variables)
    - `src` (template files to render, but this dir can be named whatever you want)
    - `bin` (template-specific scripts)

### Scripts `bin` directories and your `PATH`

DevOps-oriented shell scripts can be packaged and distributed with your jen templates and can be
either "shared" (all projects can use them, regardless of which template they use) or "template-
specific" (only accessible when a specific template is used).

When executing any action or shell command, jen always prepends your `PATH` env var with the
template-specific `bin` directory, followed by the shared one. That means you can override shared
scripts at the template level by redefining scripts with the same name as shared ones.

## Set Jen env vars

- JEN_CLONE: Local directory where jen will clone your jen git repo (defaults to `~/.jen/repo`)
- JEN_REPO: URL of your templates git repo to clone
- JEN_SUBDIR: Optional sub-directory within your repo where to look for jen files. This can be
  useful when your git repo also contains other things.

### Example

To use jen's example templates and scripts, simply set your variables as follows:

```bash
$ export JEN_REPO=git@github.com:Samasource/jen.git
$ export JEN_SUBDIR=examples
```

## Creating a new project from a template

```bash
$ mkdir foobar
$ cd foobar
$ jen do create
```

1. If it doesn't already exist, jen will automatically clone your templates git repo into your `JEN_HOME`.
2. It will prompt you to select a template from the list of those available in that repo. That template name
   will be stored in `jen.yaml` file in current dir. That file now identifies your project as jen-enabled
   and allows jen to determine the root directory of your project.
3. Depending on the steps configured in the `create` action of the template you selected, you should
   typically be prompted for the values of different variables, which will also be saved in the `jen.yaml`
   file.
4. Typically, the selected template will be rendered to current directory and, potentially, some scripts
   may automatically run to register your project with your infrastructure.

# Jen commands

## Invoke an action

To invoke any action defined in your project's template spec:

```bash
$ jen do ACTION
```

All steps defined as children of that action will be called in order.

## Execute a shell command

To execute any shell command, including your custom shell scripts, while injecting your project's env vars:

```bash
$ jen exec COMMAND ARG1 ARG2 ...
```

## Start a sub-shell

To start a sub-shell with your custom shell scripts added to `$PATH` and your project's variables as environment:

```bash
$ jen exec SHELL
```

(where `SHELL` can any of `bash`, `zsh`, `sh`...)

You are then free to call as many shell scripts and shell commands as you want, until you do `exit`.

## Update templates repo

To pull latest version of templates git repo:

```bash
$ jen pull
```

# `spec.yaml` files

Each template has a `spec.yaml` file in its root that specifies how to render the template, what
variables to prompt user and what actions user can invoke throughout the project's life-time.

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
- `exec`: executes a shell command, including shell scripts, with project vars in environment
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
    # The "exec" step allows to invoke shell commands, including shell scripts, while
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

Jen leverages the Go templating engine described [here](https://golang.org/pkg/text/template/) and
augments its built-in functions with the very helpful [sprig](https://masterminds.github.io/sprig/)
function library.

Those template expressions can be used in templates, user prompts, and file/directory names, as
described in following sections.

## Activating/deactivating rendering

By default, all files in a template are copied as is, without rendering their content as templates.
Template rendering can however be activated or deactivate selectively on a per-file/directory
basis, by appending a `.tmpl` or `.notmpl` extension to file/directory names. Applying those
extensions to a directory affects all child files recursively, unless overriden down the tree.

Note that the `.tmpl` and `.notmpl` extensions are automatically stripped away from target file/
directory names.

To override the no templating default, you can simply append a `.tmpl` extension to the name of
the root directory passed to the `render` step, ie:

```yaml
- render: ./src.tmpl
```

## Escaping double-braces

Sometimes, it's not enough to completely turn rendering on or off for an entire file. For instance,
if you need to intermix jen templating expressions with other templating that also use double-braces
(ie: helm charts) within the same file, you can escape your double-braces by using `{{{` and `}}}`,
which will be rendered to `{{` and `}}` respectively.

## Dynamic file and directory names

File and directory names can include template expressions enclosed between double-braces (ie:
`{{.PROJECT}}.sql`)

## Conditional files and directories

Files and directories can be selectively included/excluded by embedding a double-square-bracket expression
in their name, which must evaluate to true in order for the file/directory to be included in render.

Take the following template directory structure as example:

- `src`
  - `database[[.DB]]`
    - `migration.go`
    - `driver.go`

Only when the `DB` var evaluates to `true` will the `database[[.DB]]` directory and its content be rendered
to project directory. The double-square-bracket expression will also automatically get stripped away from
the target dir name:

- `src`
  - `database`
    - `migration.go`
    - `driver.go`

## Collapsing of pure conditional directories

Pure conditional directories - that is, those for which the name only contains a double-square-bracket
expression - are treated as a special case. If their expression evaluates to `true`, they get collapsed and
their contents get placed directly into parent directory.

That is very useful to group multiple files and folders under a same conditional expression, without actually
introducing an extra directory level in final output. For example, given this template structure:

- `src`
  - `[[.DB]]`
    - `migration.go`
    - `driver.go`

If `DB` is true, the following structure will be rendered to target directory:

- `src`
  - `migration.go`
  - `driver.go`

## Expressions in prompts

For prompt steps (`input`, `choice`, `option`, `options`), you can use template expressions within messages,
proposed choices and default values, by enclosing those expressions between `{{` and `}}`.

## Expressions in `if` step

As the conditional for `if` steps is always a template expression, _do not_ enclose them between double-braces, ie:

```yaml
- if: .INSTALL
  then:
    - ...
```

## Special placeholders

Because the PROJECT variable is typically used pervasively throughout templates in the form of `{{.PROJECT}}`
and `{{.PROJECT | upper}}`, we have introduced the special placeholders `projekt` and `PROJEKT`, which can
be used anywhere in file/dir names and templates without any adornments.

For example, the text "MY PROJEKT FILE.TXT" is equivalent to "MY {{.PROJECT | upper}} FILE.TXT".

Currently, those two placeholders are hardcoded and are the only ones supported, but we plan to add support
for defining your own in the template spec.

This feature was inspired by the way we were previously creating new projects by duplicating an existing
project and doing a search-and-replace for the project name in different case variants. That strategy was
very simple and effective, as long as the project name was a very distinct string that did not appear in
any other undesired contexts, hence our choice of `projekt` as something that you are (hopefully!) very
unlikely to encounter in your project for any other reason than those placeholders!

## Wishlist

- Add config to specify templates sub-dir within git repo.
- Add `confirm` step (similar to `if`, but `confirm` property contains message to display and `then` the steps to execute).
- Add `jen export` command to output env variables in a format that can be sourced directly.
- Allow `do` step to define multiple actions to call.
- Invoking `jen do` without specifying an action should prompt user to select it from available list of actions.
- Add reusable modules (including both templates and scripts).
- Add `set` step to set multiple variables.
- Add `--dry-run` flag (automatically turns on `--verbose`?).
- Add regex validation for `input` prompt.
- Allow to customize placeholders in spec file:

```
placeholders:
  projekt: {{.PROJECT | lower}}
  PROJEKT: {{.PROJECT | upper}},
```
