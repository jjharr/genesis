Credits
======

Emerald is based heavily on other contributors code.

1. Core framework : Based on Echo with some significant changes :
- binding
- validation
- server launch
- ssl configuration
- convenience methods
- middleware

2. Validation : Based on Ozzo with changes :
- return values and context
- naming (interface and Error function)
- DB validators

3. Database :

4. Migrations :

Repositories
======

Emerald keeps its own repositories for contributor packages. The primary reason we do
this is so that we have the ability to make changes as we need them. We may need changes
in packages to integrate with other parts of Emerald, design philosophy, not general
purpose, etc. We also keep these repos version in lockstep with the core emerald version
so that it's easy to version your project to a specific version of Emerald.