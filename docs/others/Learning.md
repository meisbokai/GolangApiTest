# Document learning journey while implementing the project

## Database packages in GoLang

Commonly found packages found in online tutorials and examples

- database/sql
- GORM
- sqlx
- sqlc


[source](https://blog.jetbrains.com/go/2023/04/27/comparing-db-packages/)

## Why module path naming has v1,v2 etc...

This is a module path naming convention and it applies to major versions higher than v1.

https://go.dev/ref/mod#module-path

If the module is released at major version 2 or higher, the module path MUST end with a major version suffix like /v2. This may or may not be part of the subdirectory name. For example, the module with path golang.org/x/repo/sub/v2 could be in the /sub or /sub/v2 subdirectory of the repository golang.org/x/repo.
https://go.dev/ref/mod#major-version-suffixes

Starting with major version 2, module paths MUST have a major version suffix like /v2 that matches the major version. For example, if a module has the path example.com/mod at v1.0.0, it must have the path example.com/mod/v2 at version v2.0.0.

And echo's v4 module path can be found here.

# Domain Driven Design

Domain-Driven Design is a way of structuring and modeling the software after the Domain it belongs to. What this means is that a domain first has to be considered for the software that is written. The domain is the topic or problem that the software intends to work on. The software should be written to reflect the domain. [source](https://programmingpercy.tech/blog/how-to-domain-driven-design-ddd-golang/)