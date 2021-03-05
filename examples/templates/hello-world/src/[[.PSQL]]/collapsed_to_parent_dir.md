# Collapsing to parent directory

This whole directory will be rendered only if the `PSQL` variable is `true`. Furthermore, because its name only contains the double-bracket expression `[[.PSQL]]` (and would therefore be left with an empty string, once the expression is trimmed away from final directory name), that whole directory's contents will be collapsed into parent directory.

That can be useful to group out many conditional files, without having to tag them individually with same conditional expression.
