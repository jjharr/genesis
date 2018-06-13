
Emerald Architecture Philosophy
======

These are the principles used to build Emerald. Contributors should try to follow
the same guidelines. These are listed in decreasing order of importance.

1. Maintainability. Go is designed to be easy to read and maintain. Emerald further
emphasizes this by avoid magic (init(), package-level variables, some kinds of tags, 
unnecessary reflection), and trying to localize relevant code to the places where it 
is used. (easy to debug with helpful logs and error messages). commercial support,

2. Security. Emerald aims to be secure by default while also making it clear what the
framework is doing to prevent frustration in cases where security produces unexpected
behavior. Emerald supports LetsEncrypt out of the box, with a robust SSL configuration.
It strips potentially dangerous tags from all requests except HEAD or GET by default. It
sets secure headers, uses a more robust HTTP/s server, CSRF, secure session, ...

3. Productivity. Emerald aspires to help developers get $#!* done fast. So
Emerald tries to make it easy to incorporate functionality that is very widely used in
web applications : validation, logging, database access, authorization, 
roles and permissions, transactional email, payment methods, analytics, caching, and more.
Good documentation. inversion of control. Lots of examples. Working app right
out of the box.

4. Reliability. Tests for almost all functionality. Tests must always validate relevant edge
cases such as nil pointers, zero values, very large values, multi-byte text and so on. More
robust server, try to fail at compile or launch time instead of run-time.

5. Flexibility. Almost every application has special needs.
Composable interfaces for almost everything : Binder, validation,
email, logging, database, roles, etc. Makes it possible for you to modify, use other
packages, or roll your own without having to modify the core framework. Relatively
easy to split up microservices. Use as one repo or multiple. Use for consumers or
businesses.

6. Performance. While performance is less important than the first four objectives, Go
makes it easier than most languages to avoid performance tradeoffs. Emerald enhances
performance by trying to avoid unnecessary reflection or parsing, extra database traffic, and by
using Go routines for tasks like transactional email. First class support for websockets,

It's challenging to balance these priorities. Please feel free to help us make it better!