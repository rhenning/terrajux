[![current-release](https://img.shields.io/github/release/rhenning/terrajux.svg?style=for-the-badge)](https://github.com/rhenning/terrajux/releases/latest)
[![build-status](https://img.shields.io/github/workflow/status/rhenning/terrajux/test/main?style=for-the-badge)](https://github.com/rhenning/terrajux/actions/workflows/test.yml?query=workflow%3Atest+branch%3Amain)
[![report-card](https://goreportcard.com/badge/github.com/rhenning/terrajux?style=for-the-badge)](https://goreportcard.com/report/github.com/rhenning/terrajux)
[![license-info](https://img.shields.io/github/license/rhenning/terrajux?style=for-the-badge&color=orange)](https://www.apache.org/licenses/LICENSE-2.0)

[![pb-goreleaser](https://img.shields.io/badge/%20%20%F0%9F%93%A6%F0%9F%9A%80-goreleaser-hotpink.svg?style=for-the-badge)](https://github.com/goreleaser)
[![pb-semanticrelease](https://img.shields.io/badge/%20%20%F0%9F%93%A6%F0%9F%9A%80-semantic--release-hotpink.svg?style=for-the-badge)](https://github.com/semantic-release/semantic-release)
[![pb-golangcilint](https://img.shields.io/badge/%20%20%F0%9F%93%A6%F0%9F%9A%80-golangci--lint-hotpink.svg?style=for-the-badge)](https://golangci-lint.run/)
[![pb-gosec](https://img.shields.io/badge/%20%20%F0%9F%93%A6%F0%9F%9A%80-gosec-hotpink.svg?style=for-the-badge)](https://securego.io/)


<!--
[![Codecov branch](https://img.shields.io/codecov/c/github/rhenning/terrajux/main.svg?style=for-the-badge)](https://codecov.io/gh/rhenning/terrajux)
[![Go Doc](https://img.shields.io/badge/godoc-reference-blue.svg?style=for-the-badge)](http://godoc.org/github.com/rhenning/terrajux)
-->

---

# terrajux

`terrajux` [diff](https://en.wikipedia.org/wiki/Diff)s the source code of a
[terraform](https://github.com/hashicorp/terraform)
[root module](https://www.terraform.io/docs/language/modules/index.html#the-root-module) project,
along with the source of all its transitive module dependencies, between two git refs.


## how?

- ensure that `terraform` and `git` are installed and available in your system's `PATH`.
- download the [latest release archive](https://github.com/rhenning/terrajux/releases/)
  named for the os and arch appropriate for your system.
- decompress and extract the archive.
  - `tar -zxvf terrajux_<version>_<os>_<arch>.tar.gz` will do on many systems.
- `terrajux` is distributed as a self-contained binary. simply move the `terrajux`
  binary to a location of your choosing. for many, this will be somewhere in your
  system's `PATH` such as `/usr/local/bin/` or `~/bin/`.
- on macos it may be necessary to explictly inform the system that `terrajux` is
  an approved program.
  - navigate to the directory containing `terrajux` via terminal or finder.
  - run `open .` if using the terminal.
  - right-click on `terrajux` and select **open**.

**try it!**

```
terrajux giturl v1ref v2ref [subpath]

# example:
terrajux https://github.com/terraform-aws-modules/terraform-aws-iam.git v3.15.0 master modules/iam-user
```


## why?

terraform provides primitives for building and maintaining complex, distributed
infrastructure projects through [composition of reusable modules](https://www.terraform.io/docs/language/modules/develop/composition.html). module dependencies can be [stored in source control systems](https://www.terraform.io/docs/language/modules/sources.html),
[versioned, pinned](https://www.terraform.io/docs/language/modules/sources.html#selecting-a-revision),
and [automatically fetched at runtime by `terraform init`](https://www.terraform.io/docs/cli/commands/init.html#child-module-installation). these tools provide infrastructure engineers the means to
consistently reproduce infrastructure delivery tooling with the exact version of all module and
provider dependencies used previously. they may then choose to try a new version of some module in
the development environment, and only promote those changes to staging or production when it is safe
to do so.

such a system has many benefits but also some tradeoffs. a complex infrastructure project may have
tens of modules, and each of those modules may have their own module dependencies, and so forth.

during the course of troubleshooting, outage investigation, or postmortems it can be helpful see
exactly what code changed across two versions of a project's root module and all of its dependencies
without bouncing between artifacts and source code distributed among many git repositories.

**`terrajux` aims to provide a user-friendly tool for viewing such changes.**

`terrajux` does not embed terraform. i recommend using the excellent
[`tfenv` version manager](https://github.com/tfutils/tfenv) to manage multiple
terraform versions and maintaining appropriate `.terraform-version` files in your
root projects and submodules to limit side effects of version incompatibilities.
have you ever inadvertently upgraded a project's shared terraform remote state by
using the wrong version of terraform locally? i have. don't do that.


## who?

`terrajux` is primarily aimed at site reliability and infrastructure engineers managing systems
built from nested terraform modules spread across many repositories.

it is maintained primarily by [rich henning](https://github.com/rhenning), a software engineer living and
working in philadelphia.

[@mdb](https://github.com/mdb) also contributes to this project and created [terrajux-action](https://github.com/mdb/terrajux-action), which facilitates the use of `terrajux` in a [GitHub Actions](https://github.com/features/actions) workflow.


## why "terrajux"?

**why not `terradiff`?**

[`terradiff`](https://github.com/jml/terradiff) is an existing project used to detect drift of
terraform-managed infrastructure.

**why not `________`?**

`tdiff` is concise but used as the name of some command line utilities and tree difference
libraries.

also, many of the tools intended for use within the terraform ecosystem have names beginning with
`terra`, so i wanted to stick with that convention.

`terrac[o]mp` translates from speech to text with some ambiguity.

_juxtapose_ popped into my head while considering the possibilities, and so we have `terrajux`.


## faq

> what is happening behind the scenes?

in short, `terrajux`:
- performs shallow clones of the the specified git repository at
  the supplied refs
- initializes terraform with `-backend=false` to download all module
  dependencies without touching the terraform state backend
- displays a recursive diff of the initialized projects


> what platforms are supported?

- macos and linux builds are fully tested for every pr and release
- *bsd and solaris builds are untested but _should_ work
- windows builds are disabled pending [some portability issues](https://github.com/rhenning/terrajux/issues)


> how can i use ________ to view the diff?

try the `-difftool` option.

`-difftool` accepts a [go template](https://golang.org/pkg/text/template/) that
can be used to format commands for an alternative diff viewer. the strings
`{{.V1}}` and `{{.V2}}` will be replaced by the `v1ref` and `v2ref` args
supplied to `terrajux` at runtime.

to use the [compare folders plugin](https://marketplace.visualstudio.com/items?itemName=moshfeu.compare-folders)
in [vs code](https://code.visualstudio.com/), for example, try:

```
-difftool 'COMPARE_FOLDERS=DIFF code {{.V1}} {{.V2}}'
```

to avoid typing this every time, consider creating an alias in
[your shell's profile](https://en.wikipedia.org/wiki/Unix_shell#Configuration_files),
such as:

```
alias terrajux="terrajux -difftool 'opendiff {{.V1}} {{.V2}}'"
```

> is it possible to use this in GitHub's pull request workflow?

sure! please check out [@mdb](https://github.com/mdb)'s [terrajux-action](github.com/mdb/terrajux-action) for [GitHub Actions](https://github.com/features/actions).

> i'm seeing a stale diff for a branch ref or getting strange errors during
  initialization. what gives?

try clearing the cache with the `-clearcache` option.

`terrajux` keeps a local cache of checkouts and initialized terraform modules
to speed up subsequent diffs of long-lived release tags. the cache key is a
concatentation of git url (sans scheme) and ref. when diffing dynamic refs
such as branches (or a tag that has been deleted and repointed), the cache entry
may be stale. you'll know this is the case if output displays something like
`Found <repo>@<ref> in cache. Skipping clone.`


> how do i report a bug?

woops, sorry about that! please [submit a bug report](https://github.com/rhenning/terrajux/issues/new/choose).


> how do i request a feature?

please [submit a feature request](https://github.com/rhenning/terrajux/issues/new/choose) to start the conversation. once that's done, we can work on it or
consider [pull requests](https://github.com/rhenning/terrajux/pulls).


> how can i contribute (documentation, code, resources, etc.)?

awesome! your kindness is very much appreciated. please
[check out the contribution guidelines](https://github.com/rhenning/terrajux/contribute)
for more information.


## license

this project is released under the [apache 2.0 license](LICENSE).
