# Learning Process and Challenges

## JWT Authentication

Using middlewares to implement jwt authentication was much simpler than expected because, as the name suggests, the jwt authentication simply acts as a middleware between the router and the handler. Beyond the handler, nothing needed to be changed. 

However, I had to refer to many different examples found online as there were several different approaches to attaching jwt to the API endpoints. I found that using a middleware was the cleanest and easiest to understand.

Refering to the many boilerplate codes online, other middlesware components such as OAuth, CORS and caching (Redis) have similar patterns, where the only codes affected are the ones that are layered adjacent to the middleware. However, it is unclear if this pattern is only possible due to the Clean Architecture by Robert Martin (Uncle Bob).

I was not able to completely follow the entire pattern due to the large number of foreign topics. However, I thoroughly enjoyed the process of learning many things just from implementing this project. 

# Other learning points

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

## Domain Driven Design

Domain-Driven Design is a way of structuring and modeling the software after the Domain it belongs to. What this means is that a domain first has to be considered for the software that is written. The domain is the topic or problem that the software intends to work on. The software should be written to reflect the domain. [source](https://programmingpercy.tech/blog/how-to-domain-driven-design-ddd-golang/)

## What is `v ...interface{}` as a variable (Found in log.logrus.go)

The signature of Printf uses the type ...interface{} for its final argument to specify that an arbitrary number of parameters (of arbitrary type) can appear after the format. [source](https://go.dev/doc/effective_go)

A type implements any interface comprising any subset of its methods and may therefore implement several distinct interfaces. For instance, all types implement the empty interface:
```
interface{}
```
This basically means that any type an be represented as the "empty interface", and thus that Printf can accept variables of any type for these varargs. [source](https://stackoverflow.com/questions/18629294/what-is-interface-as-an-argument)

## Testify and Mockery

Running Mockery to generate mocks for testing purposes

```
docker run -v "$PWD":/src -w /src vektra/mockery --all
```

