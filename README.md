# Jen

> Jen - _noun_ (in Chinese philosophy) a compassionate love for humanity or for the world as a whole.

Jen is a CLI tool for scaffolding new microservices based on Go templates, onboarding them with your CI/CD and infra, and augmenting them with your DevOps scripts for their entire life-time.

# Scaffolding, rendering, templating...

Throughout this document, the terms "scaffolding", "rendering" and "templating" are used interchangeably and all basically refer to the same idea of creating a project's general skeleton and boilerplate code from templates files.

# Motivation

We were not satisfied with existing project scaffolding tools, which often are language-specific and leave you on your own once your project has been generated. Many DevOps shell scripts (ie: registering your project with CI/CD and other infra, promoting it from staging to prod...) need to be maintained, organized and shared separately. They typically require similar inputs to those provided during scaffolding (ie: project name, cluster, cloud region, team...), yet you have to pass that information as arguments all over again every time you invoke them.

As DevOps, we have many concerns to address, such as:

- How do we organize and share reusable project templates?
- What about project-specific and common DevOps shell scripts?
- How do we deal with those scripts' arguments/variables?

Jen aims to provide a very simple framework to answer all those questions.

# How it works

- Put all you project templates in a git repo with a specific structure, including shell scripts specific to each template or common to all.
- Jen automatically clones that repo locally upon first use.
- When scaffolding a project, Jen prompts user for which template to use and all values required by template and associated shell scripts.
- All that information gets stored in a `jen.yaml` file in the project's root dir and remains available everytime you run a `jen ...` command anywhere within your project's directory structure.
- Use jen throughout your project's life-time to run companion shell scripts using same variables.

# Getting started

## Installation

### Installing via brew on MacOS (recommended)

```bash
$ brew tap silphid/jen
$ brew install jen
```

### Downloading binary

- Download and install latest release for your platform from the GitHub [releases](https://github.com/silphid/jen/releases) page.
- Make sure the binary is accessible via your `$PATH`.

### Building and installing from source

This approach requires that you replace both mentions of version number with your desired version in the following command:

```bash
$ go install -ldflags "-s -w -X 'main.version=v0.1.2'" github.com/silphid/jen/cmd/jen@v0.1.2
```

## Create git repo for templates/scripts

Create a git repo to store all your templates and scripts, using the following structure:

- `bin` (scripts common to all templates)
- `templates`
  - `TEMPLATE_NAME`
    - `spec.yaml` (defines actions/steps/variables)
    - `project` (template files to render, but this dir can be named whatever you want)
    - `bin` (template-specific scripts)

### Scripts `bin` directories and your `PATH`

DevOps-oriented shell scripts that rely on your project variables can be packaged and distributed with your jen templates. Those scripts can be placed in three kinds of locations (in order of precedence):

- project-level `bin` dir in project's root contains project-specific scripts (only accessible from that project).
- template-level `bin` dir contains template-specific scripts (only accessible from projects associated with that template).
- top-level `bin` dir in git repo contains shared scripts (accessible from all projects regardless of which template they are associated with).

When executing any action or shell command, jen always prepends your `PATH` env var with the three directories above, in that order. That means you can override scripts at a more specific level by redefining them with the same names as less specific ones.

## Set Jen environment variables

- `JEN_CLONE`: Local directory where jen will clone your jen git repo (defaults to `~/.jen/repo`)
- `JEN_REPO`: URL of your templates git repo to clone
- `JEN_SUBDIR`: Optional sub-directory within your repo where to look for jen files. This can be useful when your git repo also contains other things.

# Hello World example

This `hello-world` example is available on [github](https://github.com/silphid/jen/tree/master/examples/templates/hello-world). Don't hesitate to explore it there to better understand how it works (the template itself is much more instructive and interesting than the final output).

## Configuring jen

1. Configure jen to point to jen's example templates and scripts:

```bash
$ export JEN_REPO=git@github.com:silphid/jen.git
$ export JEN_SUBDIR=examples
```

## Creating project

2. Create a new project directory:

```bash
$ mkdir foobar
$ cd foobar
$ jen do create
```

3. If it's your first time, jen will automatically clone your templates git repo into `$JEN_HOME/repo`.
4. Because the current dir is not initialized with jen yet, it asks for confirmation. Type `y` and press `Enter`.
5. Jen then shows a list of available templates from that repo. Right now there's only one `hello-world` example, so just press `Enter`. That choice gets saved to `jen.yaml` file in current dir and identifies your project as jen-initialized.
6. Because the `create` action calls out to the `prompt` action, you are now prompted for variable values. Answer the different prompts (notice how it automatically suggests the current dir name `foobar` as default project name). Your values also get saved to `jen.yaml` file.
7. The `create` action then calls `render` step to render the `hello-world` template files to current dir.
8. If in previous prompts you opted for installing your project in CI/CD, the `install` action will be called now to simulate that.
9. At this point, typically, you would commit your project to git, including the `jen.yaml` file.

## Inspecting project variables

Let's have a look at our project's `jen.yaml` file that has just been created:

```bash
cat jen.yaml
version: 2021.04
template: hello-world
vars:
  INSTALL: true
  NEWRELIC: true
  PROJECT: foobar
  PSQL: true
  TEAM: devops
```

But there's a dedicated command for viewing variables, which can be invoked from anywhere within your project structure:

```bash
$ jen list vars
INSTALL: true
NEWRELIC: true
PROJECT: foobar
PSQL: true
TEAM: devops
```

## Invoking actions

We are now ready to call different project actions with `jen do ACTION`, but first let's see what actions the `hello-world` example defines:

```bash
$ jen list actions
create
install
prompt
uninstall
```

We have already discussed about `create` and `prompt`. Now, `install` and `uninstall` are meant to register/unregister your project with your CI/CD pipeline and infra, but here they just call dummy bash scripts that simulate the real thing. For example:

```bash
$ jen do install
Creating docker image repo for project foobar
Done.
Creating triggers on CI/CD pipelines for project foobar
Done.
```

You can also run `jen do` without specifying action and it will prompt you for which one to execute:

```bash
$ jen do
? Select action to execute  [Use arrows to move, type to filter]
> create
  install
  prompt
  uninstall
```

## Executing scripts

Typically, we would always go through higher-level actions to call scripts and shell commands, but we can also invoke them directly using `jen exec CMD ARG1 ARG2 ...`. However, let's first see what scripts the `hello-world` example defines:

```bash
$ jen list scripts
create-cicd-triggers
create-docker-repo
remove-cicd-triggers
remove-docker-repo
```

These are the scripts that are present either in the templates' shared `bin` dir or in this template's specific `bin` dir, if any.

Let's try one of them:

```bash
$ jen exec remove-cicd-triggers
Removing triggers from CI/CD pipelines for project foobar
Done.
```

Or even simpler, just run `jen exec` alone to let it prompt you for custom script to execute:

```bash
$ jen exec
? Select script to execute  [Use arrows to move, type to filter]
> create-cicd-triggers
  create-docker-repo
  remove-cicd-triggers
  remove-docker-repo
```

Just keep in mind that you are not limited to custom scripts, you can execute really any shell command with the `jen exec CMD ARG1 ARG2 ...` syntax.

## Starting a sub-shell

You can even start a sub-shell with your custom shell scripts added to `$PATH` and your project's variables as environment:

```bash
$ jen shell
```

You are then free to call as many shell scripts and shell commands as you want, until you do `exit`. For example:

```bash
$ echo $PROJECT
foobar

$ remove-cicd-triggers
Removing triggers from CI/CD pipelines for project foobar
Done.

$ exit
```

Note that this command is just a shorthand for:

```bash
$ jen exec $SHELL
```

Where the `$SHELL` variable is typically set to your current shell, but you can also explicitly specify any of `bash`, `zsh`, `sh`...

## Updating templates git repo

To pull latest version of templates git repo:

```bash
$ jen pull
```

## Cleaning up

When you're done experimenting with the examples, don't forget to delete the jen examples repo clone from your machine (at `$JEN_HOME/repo` or `~/.jen/repo`) and to make jen point to your own template repo.

# `spec.yaml` files

Each template has a `spec.yaml` file in its root that specifies how to render the template, what variables to prompt user and what actions user can invoke throughout the project's life-time.

It has this general structure:

```yaml
version: 2021.04
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

Actions are named operations that can be invoked by user via the `jen do ACTION` command, or as part of another action (using the `do` step). The order of actions is irrelevant, much like the order of function definitions within any source code.

### Standard actions

You can have any arbitrary actions with any names in your template specs, however it is recommended to follow the convention of having at least the following actions:

- `create`:
  - first invoke the `prompt` action below
  - then render project template
- `prompt`:
  - prompt user for variable values

Optionally, also include those actions:

- `install`:
  - register your project with CI/CD pipelines and infra
- `uninstall`:
  - unregister your project from CI/CD pipelines and infra

## Steps

Each action is composed of one or many steps that are executed sequentially when the action is invoked (their order is therefore important).

Steps have predefined names and purposes:

- `if`: conditionally invokes child steps
- `do`: executes another action by name (much like a function call)
- `exec`: executes a shell command, including shell scripts, with project vars in environment
- `render`: renders template into current dir, using project vars
- `input`: prompts user for a single free-form string var
- `choice`: prompts user for a single string var among a list of multiple proposed choices
- `option`: prompts user for single boolean var as a yes/no question
- `options`: prompts user for multiple boolean vars as a list of toggles
- `set`: sets one or multiple variables to given values without user intervention

## Example

The following is the [spec.yaml](https://github.com/silphid/jen/tree/master/examples/templates/hello-world/spec.yaml) file of `hello-world` template in jen's [examples](https://github.com/silphid/jen/tree/master/examples):

```yaml
# Version of jen file format (for future compatibility checks)
version: 2021.04

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
        # Here we use the special "PROJECT_DIR_NAME" variable to propose to the user the
        # project's directory name as default project name, which is most often the case.
        default: "{{ .PROJECT_DIR_NAME }}"

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

    # The "set" step sets one or multiple variables to given values without user intervention.
    - set:
        ORGANIZATION_ID: 123456789
        CLOUD_PROJECT_ID: 987654321

  # By convention, the "install" action is in charge of setting the project up with CI/CD
  # and infra. It is typically invoked as the last step of the "create" action.
  install:
    # The "exec" step allows to invoke shell commands, including shell scripts, while
    # passing them all project variables as env vars. In this case, it specifies a list
    # of multiple commands to execute.
    - exec:
        - create-docker-repo
        - create-cicd-triggers

  # By convention, the "uninstall" action is in charge of removing the project from infra
  uninstall:
    # The "confirm" step is similar to "if", however it prompts user with given message and
    # only upon confirmation executes steps in the "then" clause.
    - confirm: Are you sure you want to completely uninstall project {{.PROJECT}} from infrastructure?
      then:
        # Here the "exec" step is invoked multiple times, each executing a single command
        - exec: remove-docker-repo
        - exec: remove-cicd-triggers

  # This action can be invoked multiple times after project has been initially scaffolded, in
  # order to simulate adding endpoints to our microservice.
  add-endpoint:
    - input:
        question: Endpoint name
        # Variables prefixed with ~ are transient (temporary) and are not saved to jen.yaml file.
        # Note that the ~ symbol is not actually part of the variable name.
        var: ~NAME
    - input:
        question: Endpoint path
        var: ~URL_PATH
    # This renders templates under `endpoint` sub-directory. One of those template files has a `.insert`
    # extension, meaning it is a special "insertion" template, which is meant to be inserted into an
    # existing file of the same name at a specific insertion point (defined by regular expressions).
    # Have a look at that file for more details on how it works.
    - render: endpoint
```

# Templates

## Go template language

Jen leverages the Go templating engine described [here](https://golang.org/pkg/text/template/) and augments its built-in functions with the very helpful [sprig](https://masterminds.github.io/sprig/) function library.

Those template expressions can be used in templates, user prompts, and file/directory names, as described in following sections.

## Activating/deactivating rendering

By default, all files in a template are copied as is, without rendering their content as templates. Template rendering can however be activated or deactivate selectively on a per-file/directory basis, by appending a `.tmpl` or `.notmpl` extension to file/directory names. Applying those extensions to a directory affects all child files recursively, unless overriden down the tree. Note that the `.tmpl` and `.notmpl` extensions are automatically stripped away from target file/directory names.

To override the no templating default, you can simply append a `.tmpl` extension to the name of the root directory passed to the `render` step, ie:

```yaml
- render: ./src.tmpl
```

## Escaping double-braces

Sometimes, it's not enough to completely turn rendering on or off for an entire file. For instance, if you need to intermix jen templating expressions with other templating that also use double-braces (ie: helm charts) within the same file, you can escape your double-braces by using `{{{` and `}}}`, which will be rendered to `{{` and `}}` respectively.

## Dynamic file and directory names

File and directory names can include template expressions enclosed between double-braces (ie: `{{.PROJECT}}.sql`)

## Conditional files and directories

Files and directories can be selectively included/excluded by embedding a double-square-bracket expression in their name, which must evaluate to true in order for the file/directory to be included in render.

Take the following template directory structure as example:

- `src`
  - `database[[.DB]]`
    - `migration.go`
    - `driver.go`

Only when the `DB` var evaluates to `true` will the `database[[.DB]]` directory and its content be rendered to project directory. The double-square-bracket expression will also automatically get stripped away from the target dir name:

- `src`
  - `database`
    - `migration.go`
    - `driver.go`

## Collapsing of pure conditional directories

Pure conditional directories - that is, those for which the name only contains a double-square-bracket expression - are treated as a special case. If their expression evaluates to `true`, they get collapsed and their contents get placed directly into parent directory.

That is very useful to group multiple files and folders under a same conditional expression, without actually introducing an extra directory level in final output. For example, given this template structure:

- `src`
  - `[[.DB]]`
    - `migration.go`
    - `driver.go`

If `DB` is true, the following structure will be rendered to target directory:

- `src`
  - `migration.go`
  - `driver.go`

## Expressions in prompts

For prompt steps (`input`, `choice`, `option`, `options`), you can use template expressions within messages, proposed choices and default values, by enclosing those expressions between `{{` and `}}`.

## Expressions in `if` step

As the conditional for `if` steps is always a template expression, _do not_ enclose them between double-braces, ie:

```yaml
- if: .INSTALL
  then:
    - ...
```

## Special placeholders

Placeholders are a lightweight alternative to go template expressions, which can be used as plain text anywhere in file/dir names and template files. Because placeholders are processed using plain search-and-replace, ensure they have improbable names that don't risk conflicting with anything else (ie: "projekt").

For example, you can define the following placeholders in your template spec:

```yaml
placeholders:
  projekt: "{{ .PROJECT | lower }}"
  Projekt: "{{ .PROJECT | title }}"
  PROJEKT: "{{ .PROJECT | upper }}"
```

You can then use these placeholders anywhere without any adornments. For example, the text "MY PROJEKT FILE.TXT" is equivalent to "MY {{.PROJECT | upper}} FILE.TXT".

This feature was inspired by the way we were previously creating new projects by duplicating an existing project and doing a search-and-replace for the project name in different case variants. That strategy was very simple and effective, as long as the project name was a very distinct string that did not appear in any other undesired contexts, hence our choice of `projekt` as something that you are (hopefully!) very unlikely to encounter in your project for any other reason than those placeholders!

## Adding multiple similar elements to a project after scaffolding

Let's say you want developers to be able to add multiple endpoints to a microservice, each one with its own sub-dir and source files. To achieve that you simply need to put your endpoint template files in a separate sub-dir than the main template files. For example, if your project's main template files are in a `project` sub-dir, you could create another `endpoint` sub-dir with just your endpoint template files. Then simply create a standalone action that prompts user for endpoint-specific values and then renders the `endpoint` sub-dir using those values.

See `hello-world` example template for a demonstration of adding multiple endpoints to an already generated project.

## Inserting content into an existing file at a given location

The endpoint scenario described in previous section is fine, except that the files and directories you generate for each endpoint will typically not just sit there in your project. You probably also need to reference them from some parent source file. That means that for each endpoint you add to the project, you would need to insert referencing code into some existing file.

To that end, Jen supports special template files named "inserts" and marked with the `.insert` extension (or with `.insert.` anywhere in their name) that are intended to be inserted into a target file of same name (minus the `.insert` extension) that must already exist in project at same path location.

Each insert template file may define one or more insertion sections, each delimited by `<<< START_REGEX` and `>>> END_REGEX` lines. The `START_REGEX` and `END_REGEX` expressions are optional, but at least one of them must be specified. Those start and end regular expressions allow to find the insertion location in target file for the section's template body. For example:

```
<<< ^List of endpoints
Definition of endpoint {{.NAME}} for path {{.PATH}}
>>> ^$
```

The start regex above (`^List of endpoints`) serves to find first line that starts with `List of endpoints`, then the end regex (`^$`) serves to find first empty line following start line. Jen will then insert the template's body (`Definition of endpoint...`) before that empty line.

If you need to insert text in different locations of same file, you can specify multiple sections, each delimited by `<<<` and `>>>` markers.

All text outside delimited sections simply gets discarded/ignored.

The rules for determining insertion point are the following:

- If you specify only start regex, insertion will happen right after first matching start line.
- If you specify only end regex, insertion will happen right before first matching end line.
- If you specify both start and end regexes, insertion will happen right before first matching end line that is found following first matching start line (in other words, end line is searched relatively to start line).

See `hello-world` example template for a demonstration of inserting multiple snippets into an existing source file at a specific insertion location.

For complete regex syntax reference, see the [RE2 wiki](https://github.com/google/re2/wiki/Syntax).

# Other commands

## Verifying required variables in custom scripts

To make your jen bash scripts more robust and self documented, you can use the `jen require VAR1 VAR2 ...` command in their first few lines (typically after `set -e` to make script fail in case of missing variable):

```bash
#!/bin/bash
set -e
jen require PROJECT TEAM
echo "You are now garanteed that the $PROJECT and $TEAM variables can be used safely"
```

# Tips

## Associating an existing project with a template

To associate a template with an existing project that was not initially generated by jen, without doing any scaffolding, you just have to invoke the `jen do prompt` command in the root of the existing project. This assumes your templates follow the recommended convention of having the standard `create` and `prompt` actions (where the `create` action first calls `prompt` and then does the template rendering). In that case, calling the `prompt` action alone in a non-jen-initialized project will first ask you to select the template to associate the project with, and then will prompt you for variable values and save them to the `jen.yaml` file. From that point, your project is initialized and associated with a template. You just need to commit the `jen.yaml` file into git.

# Wishlist

- Support for versioning templates, allowing projects generated with a specific template version to use a specific version of bash scripts (useful for progressively evolving scripts without breaking older projects).
- Add `jen which [script]` command to show absolute path to script that would be executed by `jen exec [script]`.
- Auto-completion for bash and zsh shells.
- Add `--dry-run` flag (automatically turns on `--verbose`?).
- Add regex validation for `input` prompt.
- Add `jen confirm MESSAGE` command for scripts to use for confirming dangerous operations like uninstalling (the command would return exit code 0 or 1, depending on whether user responds Yes or No respectively).
- Fix `choice` step to pre-select current value, if any.
- Way for overridding bin script to fallback/call-out to default implementation(?)
- Add reusable modules, including both templates and scripts (ie: scripts common to all go projects).
