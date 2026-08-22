# Git plugin

The Git MVP recognizes status, add, commit, diff, log, branch, switch, checkout, restore, stash, pull, push, fetch, remote, merge, rebase, init, and clone.

Detection runs Git directly and parses repository membership, current branch, porcelain status, local branches, remotes, upstream, and commits ahead. Core caches the returned state for one second. Dynamic completion covers local branches for switch/checkout/delete/merge/rebase, changed files for add/restore, and remotes for push/pull/fetch.

Priority metadata reflects actual state: modified files raise diff/add relevance, staged files raise cached-diff/commit relevance, and ahead commits raise push relevance. Core still owns final ranking. Successful add/commit/pull/status commands produce contextual next actions. Reviewing staged changes is a first-class best practice. Failed repository commands can suggest init, while pathspec/revision failures suggest real branches.

Limitations: porcelain rename paths receive only minimal parsing; remote branches and detailed option completion are not included; the short cache can briefly show stale state; recovery intentionally recognizes only common failure text.
