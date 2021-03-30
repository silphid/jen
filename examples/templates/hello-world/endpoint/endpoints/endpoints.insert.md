# This file defines insertions

Multiple insertion sections can be defined here, each delimited by `<<< START_REGEX` and `>>> END_REGEX` 
lines. `START_REGEX` and `END_REGEX` are optional, but at least one of them must be specified.

All text outside delimited sections gets discarded/ignored. In this case the only lines that matter
are lines 14-16, the rest is just comments.

In this file there is only a single insertion section, the next paragraph. The start regex here serves
to find first line that starts with `List of endpoints`. Then, the end regex serves to find first
empty line (denoted by `^$`) that follows start line. Jen will then insert a single line of text (that
is `Definition of endpoint...`) before that empty line.

<<< ^List of endpoints
Definition of endpoint {{.NAME}} for path {{.PATH}}
>>> ^$

The rules for finding insertion point is as follows:
- If you specify only start regex, insertion will happen right after first matching start line.
- If you specify only end regex, insertion will happen right before first matching end line.
- If you specify both start and end regexes, insertion will happen right before first matching end line after
first matching start line.
