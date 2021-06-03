# broken windows

* git
    * supports only branch & tag refs atm, consider arbitrary (short) sha1's
* testing
    * mocks & stubs - many tests run real code. use the included interfaces to
      improve tests and avoid calling code with side effects
    * abstract external command executor via its own package
* delivery
    * homebrew
* documentation
    * use concrete example based on fixtures in "try it"
    * add godoc to packages
    * named return values on all functions
    * consider additional documentation
    * contributing guidelines
* configuration
    * expose configuration overrides (difftool, terraform) to user
* errors
    * better feedback in error handling
