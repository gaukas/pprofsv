# `pprofsv`: pprof-based State-transition Valitator

`pprofsv` is a project to validate (software) state-transition models on pprof profiles with user-defined assertions in Go.

## Development Status
- [x] Support of `pprof` profile output
    - [x] Support generalized call stack
    - [ ] Support inline functions
- [x] Support of user-defined assertions
    - [x] in Go
    - ...and other languages/formats
- [x] SAT-based Model Checking
