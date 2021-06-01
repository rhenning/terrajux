# broken windows

* git
    * supports only branch & tag refs atm, consider arbitrary (short) sha1's
* testing
    * mocks & stubs - many tests run real code. use the included interfaces to
      improve tests and avoid calling code with side effects
    * abstract external command executor via its own package
* delivery
    * consider github actions after better mocks & stubs
    * consider goreleaser for build/release management
    * test, vet, fmt, static analysis checks
    * homebrew
* documentation
    * use concrete example based on fixtures in "try it"
    * add godoc to packages
    * named return values on all functions
    * consider additional documentation
    * contributing guidelines
    * github metadata files
    * FAQ
* configuration
    * expose configuration overrides (difftool, terraform) to user
* errors
    * better feedback in error handling
