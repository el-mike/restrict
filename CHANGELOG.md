# 2.0.0

- Allows Subject to have multiple roles
- Errors API changes:
  - `AccessDeniedError` now has `Reasons` property, containing one or more `PermissionErrors` describing what exactly went wrong
    for every Role of the Subject and Action in the Request
  - `PermissionError` now has `ConditionErrors` property, containing zero or more `ConditionNotSatisfiedErrors` describing `Condition` failures for given Permissions
- properties of `AccessDeniedError` and `PermissionError` are now public
- Allows to perform complete validation when `AccessRequest.CompleteValidation` is set to true
  - Defaults to false (fail-early strategy)
  - When using `CompleteValidation`, all Policy errors will be collected and returned at once, instead of failing on
    first encountered Policy error
