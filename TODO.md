# broken windows

* cli
    * missing, implement it
* git
    * supports only branch & tag refs atm, consider arbitrary (short) sha1's
* testing
    * mocks & stubs - many tests run real code. refactor to interfaces to improve
      testability and avoid calling code with side effects
* delivery
    * consider github actions after better mocks & stubs
    * consider goreleaser for build/release management
    * test, vet, fmt, static analysis checks
* documentation
    * use concrete example based on fixtures in "try it"
    * add godoc to packages
    * consider additional documentation
    * contributing guidelines
    * github metadata files
* configuration
    * expose configuration overrides (difftool, terraform) to user
