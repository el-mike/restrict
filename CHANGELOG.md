# 2.0.0

- Allows Subject to have multiple roles
- Errors API changes:
  - `AccessDeniedError` now has `Errors` property, containing `PermissionErrors` describing what exactly went wrong
    for every Role of the Subject and Action in the Request
  - `AccessDeniedError` and `PermissionError` properties are now public
- Allows to perform complete validation when settings `AccessRequest.CompleteValidation` to true
  - Defaults to false (fail-early strategy)
  - When using `CompleteValidation`, all Policy errors will be collected and returned at once, instead of failing on
    first encountered Policy error
