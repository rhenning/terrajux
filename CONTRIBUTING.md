# contributing to terrajux

thank you for considering a contribution to `terrajux`! we'd love to have you
onboard. please keep the following guidelines in mind when participating.


## code of conduct

*tl;dr `don't be a jerk.`*

i don't want to belabor the point but please try to be a decent human being
when corresponding with other human beings about the project. **open source is
a labor of love.** keep it civil, on-topic, be patient, and show empathy. don't
use derogatory or condescending language, or name things in ways that might
make others uncomfortable. be kind in issues and in code reviews. `terrajux`
has no warranty. while reasonable effort is made to be careful, software has
bugs, and we are terribly sorry if it corrupts your hard drive, drinks all
your beer, or has an accident on the carpet.


## question or problem?

feel free to [open an issue](https://github.com/rhenning/terrajux/issues/new/choose). for bugs, choose **bug report**. for feature requests or changes,
choose **suggestion**. for all other inquiries please choose the "blank issue"
type.


## found a bug?

if you find a bug, please [submit a bug report](https://github.com/rhenning/terrajux/issues/new/choose)
and consider creating a [pull request](#pull-requests).


## missing a feature?

request new features by [submitting a suggestion](https://github.com/rhenning/terrajux/issues/new/choose).
if the feature is significant in scope, or might alter existing program behavior
or architecture in an impactful way, please wait for a response before taking
on any implementation. small feature improvements can be submitted directly
via [pull request](#pull-requests), if desired.


## pull requests

please consider the following when creating pull requests:

- [search](https://github.com/rhenning/terrajux/pulls) for open or closed pull
  requests related to your idea.
- [open an issue](https://github.com/rhenning/terrajux/issues/new/choose) to
  describe the feature or bug and link it to your pr.
- [fork](https://docs.github.com/en/github/getting-started-with-github/fork-a-repo)
  the program repository.
- [clone](https://docs.github.com/en/github/creating-cloning-and-archiving-repositories/cloning-a-repository-from-github/cloning-a-repository)
  your fork of the repository.
- create a new branch for your changes:

    ```
    git checkout -b fix-branch main
    ```

- create your patch. be sure to create tests, run `make test`, and
  update relevant documentation.
- commit your changes with a short, descriptive commit message describing your
  change:
  
    ```
    git commit --all
    ```

- push your branch:

    ```
    git push origin fix-branch
    ```

- [create a pull request](https://docs.github.com/en/github/collaborating-with-pull-requests/proposing-changes-to-your-work-with-pull-requests/creating-a-pull-request-from-a-fork)
  targeting `terrajux:main`.


### changing a pull request

if changes need to be made to the pull request:

- make updates to the code
- re-run the tests
- create a fixup commit and push your branch to update the pull request:

    ```
    git commit --all --fixup
    git push
    ```


## releasing

while individual commit messages within branches can be somewhat freeform, *the
first line of **merge commit messages must follow a subset of the
[semantic release commit message format](https://github.com/semantic-release/semantic-release#commit-message-format)***.
merge commit messages are hints to the [continuous delivery system](https://github.com/rhenning/terrajux/actions),
used for both automatic triggering of [semantically-versioned](https://semver.org/) 
releases and changelog generation.

in short:

- `fix: <message>` indicates a patch or bug fix that does not alter the public
  program interface in an incompatible way. this *bumps the **patch version**
  and generates a release* upon merge to `main` and successful completion of
  tests.
- `feat: <message>` indicates a feature that adds functionality but does not
  alter the public program interface in an incompatible way. this *bumps the
  **minor version** and generates a release* upon merge to `main` and
  successful completion of tests.
- `perf: <message>` indicates a significant rewrite or breaking change in the
  program's public interface. this *bumps the **major version** and generates
  a release* upon merge to `main` and successful completion of tests.

the following utility message conventions provide useful context but do not
generate releases:

- `build: <message>` indicates a change to the build or compiler configuration.
- `ci: <message>` indicates a change to the continuous delivery configuration.
- `chore: <message>` indicates a non-release change not covered by other types.
- `docs: <message>` indicates a documentation-only change.
- `refactor: <message>` indicates a code change that does not fix a bug, add a
  feature, or change the program's public interface.
- `revert: <message>` indicates a previously merged change is being reverted.
- `style: <message>` indicates a non-breaking style-related change.
- `test: <message>` indicates the addition or correction of test cases.

when choosing a merge message prefix, consider the *most significant change*
included in the branch being merged. for example, if a pr being merged includes
a non-breaking bug fix and a new feature, then choose the `^feat:` strategy.


### recovering a failed release

the automated release process can go wrong for a few reasons, including bugs
in configuration or release tools, malformed commit messages, transient ci/cd
or network errors, etc.

release failures typically fall into one of just a few categories:

- a [test failure occurred](#test-failure) on-merge to trunk (git's `main` branch)
- on-merge tests passed but the [release commit was not tagged](#release-not-tagged)
  due to a [malformed commit message](#releasing), bug, configuration, or transient
  error
- the merge commit was tagged and pushed but [publishing release assets failed](#publish-release-failed)
  due to a bug, config, or transient error

#### test failure

if a test failure occurred, re-run the failed [github action for `main`](https://github.com/rhenning/terrajux/actions?query=branch%3Amain).
if tests now pass, the release process should resume automatically from where it
left off. if tests continue to fail, view the logs, try to get a sense of what
went wrong, attempt to reproduce the issue locally, then create a new pr from
`main` and attempt to fix the tests or the test config (`/.github/workflows/test.yml`)
in that branch. the test workflow is identical in topic branches and in trunk, so
if a build works in a topic branch, it _should_ work in `main`. if all else fails,
don't sweat it and [open an issue](https://github.com/rhenning/terrajux/issues/new/choose).


#### release not tagged

the release tag step can fail due to merge commit messages that don't match
`^(fix|feat|perf):`, errors in ci/cd workflow configuration, bugs, or transient
errors. if the release was not tagged, view the
[github actions for `main`](https://github.com/rhenning/terrajux/actions?query=branch%3Amain)
and try to discern why the [release job's](https://github.com/rhenning/terrajux/actions/workflows/release.yml) 
tag step failed. if it's due a configuration error, please submit a pr with a
fix or open an issue.

to manually tag a release commit:

```
git checkout <commit-that-should-have-been-tagged>
git tag -a vA.B.C -m "release A.B.C"
git push origin vA.B.C
```

see [the next section](#publish-release-failed) to build and publish release
artifacts after tagging.

*note:* when manually tagging releases, please be aware of whether the release
is a `fix` (patch release), `feat` (minor release), or `perf` (major release)
type and tag accordingly.


#### publish release failed

sometimes the release might be tagged but publishing of assets has failed. this
can be due to errors in configuration, bugs, or transient errors. you'll know
this is the case if the [tag shows up on the releases page](https://github.com/rhenning/terrajux/releases)
but the release contains no assets.

if there is a bug in configuration, please open a pr to fix the bug. you can
then perform a release against the failed tag with the repaired configuration:

```
git checkout vA.B.C
git checkout .goreleaser.yml origin/main
goreleaser release --rm-dist --skip-validate
```
