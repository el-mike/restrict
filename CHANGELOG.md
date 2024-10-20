# 2.0.0

- Allows Subject to have multiple roles
- Errors API changes:
  - `AccessDeniedError` now has `Errors` property, containing `PermissionErrors` describing what exactly went wrong
    for every Role of the Subject and Action in the Request
  - `AccessDeniedError` and `PermissionError` properties are now public
