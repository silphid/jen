# Jen

> Jen - _noun_ (in Chinese philosophy) a compassionate love for humanity or for the world as a whole.

Code generator and script runner.

## Wishlist

- `jen do` alone to prompt for action
- `jen export` to list env variables
- Reusable modules
- Add `set` step to set multiple variables (are those saved to `jen.yaml`?)
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
