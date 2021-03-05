# Conditional directory

This whole directory will be rendered only if the `PSQL` variable is `true`. Also, the `psql[[.PSQL]]` directory will be renamed to `psql`.

# Template rendering not activated

Because this file doesn't have the `.tmpl` extension (nor any of its parent directories), template rendering is not activated in it. Therefore this {{.PSQL}} expression is left as is.