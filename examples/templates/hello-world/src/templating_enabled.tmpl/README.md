# Template rendering activated recursively

Because containing directory has the `.tmpl` extension, template rendering gets activated for all files it contains, recursively. Therefore this {{.TEAM}} expression gets interpolated. The `.tmpl` extension also gets trimmed away from the final directory name.