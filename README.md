# terrajux

`terrajux` [diff](https://en.wikipedia.org/wiki/Diff)s the source code of a
[Terraform](https://github.com/hashicorp/terraform)
[root module](https://www.terraform.io/docs/language/modules/index.html#the-root-module) project
and all of its transitive module dependencies between two git refs.


## how?

try it!

```
terrajux giturl subpath ref1 ref2

# for example:
terrajux https://github.com/terraform-aws-modules/terraform-aws-iam.git /modules/iam-user 4.1.0 master
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
with nested terraform modules that may be spread across many repositories.

it is maintained by [rich henning](https://github.com/rhenning), a software engineer living and
working in philadelphia.


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
