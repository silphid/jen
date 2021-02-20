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

# Templates git repo

- bin (scripts common to all templates)
- templates
  - TEMPLATE_NAME
    - spec.yaml (defines actions/steps/variables)
    - src (template files to render)
    - bin (template-specific scripts)

# Spec files

Each template has a `spec.yaml` file in its root that specifies how to render the template, what
variables to prompt user for values and what actions to make available to user throughout the
project's lifetime.

A `spec.yaml` file has this general structure:

```yaml
metadata:
  name: TEMPLATE_NAME
  description: ...
  version: 0.2.0
actions:
  ACTION1:
    - STEP1: ...
    - STEP2: ...
    - STEP3: ...
  ACTION2:
    - STEP1: ...
```

Here's an example spec:

```yaml
metadata:
  name: go-service
  description: An archetypical go microservice
  version: 0.2.0
actions:
  create:
    - do: prompt
    - render: src
    - exec:
        - create-git-repo
        - create-ecr-repo
        - create-cicd-triggers
  prompt:
    - input:
        question: Project name
        var: PROJECT
        default: $(basename $(pwd))
    - options:
        question: Select desired features
        items:
          - text: PostgreSQL database
            var: PSQL
            default: true
          - text: NewRelic instrumentation
            var: NEWRELIC
            default: false
  destroy:
    - exec:
        - remove-git-repo
        - remove-ecr-repo
        - remove-cicd-triggers
```

## The `create` action

By convention, one of the actions is always named `create`. It is the one invoked initially to
prompt user for variable values and scaffold your project. The other actions are typically invoked
manually by user later on, but can also be called from the `create` action.

# Templates

# Shell scripts

When executing any jen action or shell command, jen automatically adds all the project's variables to your
shell environment and also the template-specific and common `bin` directories to your `$PATH`. That way,
your scripts can only be invoked when all proper environment variables are set.

## Wishlist

- Invoking `jen do` without specifying an action should prompt user to select it from available list of actions.
- Add `jen export` command to output env variables in a format that can be sourced directly.
- Add reusable modules (including both templates and scripts).
- Add `set` step to set multiple variables
- Add `confirm` step (similar to `if`, but `confirm` property contains message to display and `then` the steps to execute)
- Allow to customize placeholders in spec file:

```
placeholders:
  projekt: {{.PROJECT | lower}}
  PROJEKT: {{.PROJECT | upper}},
```

- Add `--dry-run` flag (automatically turns on `--verbose`?)
