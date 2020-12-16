# Jen

> Jen - *noun* (in Chinese philosophy) a compassionate love for humanity or for the world as a whole.

Code generator and script runner.

## Wishlist (if time allows)

- Only treat as templates files (or entire directories) ending with `.gotmpl` (and remove extension)
- When prompting again, reuse existing values as defaults
- Per-template/module scripts in `bin` dir, which are automatically included in `PATH`
  - `jen exec` and `jen export` should alter `PATH` to include `bin` dir(s)
- `jen do` alone to prompt for action
- `jen export` to list env variables
- Reusable modules
- Add `set` step to set multiple variables (are those saved to `jen.yaml`?)
- Override values at command-line level (`--set myValue=value`)
- Exceptionally escape within any file, using `{{{` and `}}}` to represent `{{` and `}}`
- `confirm` step (similar to `if`, but `confirm` property contains message to display and `then` the steps to execute)
- Custom placeholders:
```
placeholders:
  projekt: {{.PROJECT | lower}}
  PROJEKT: {{.PROJECT | upper}}
```
- `--dry-run` flag (automatically turns on `--verbose`?)

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
