taken from Globio - modify.

# Coding standards

This document is a work-in-process. There are many references here to [this document](https://github.com/golang/go/wiki/CodeReviewComments) 
which is recommended reading.

## Formatting

- Always run gofmt before commiting any code.
- When implementing an interface on a type, the order of interface methods on the type should
be the same as it is on the interface definition.

## Naming

Naming is extremely important to readability, and thus maintainability. Struct and function
names are especially important. Please take extra care to name variables and functions in a
highly descriptive way, so that someone who knows nothing about the structure of your code
could easily read a top-level method and understand what is likely happening in the functions
being called. If in doubt, err on the side of long names over short ones. Also :
* [Package Names](https://github.com/golang/go/wiki/CodeReviewComments#package-names)
* [Initialisms](https://github.com/golang/go/wiki/CodeReviewComments#initialisms)
* [Constants](https://github.com/golang/go/wiki/CodeReviewComments#mixed-caps)
* Functions that act on the DB :
    * Insert : InsertX<ObjectName>
    * Select : By<PropertyName>
    * Update with many fields : Update<ObjectName>
    * Update of one or more related flag fields : MarkAs<FlagName>


## Function Arguments

- [Pass a struct by value](https://github.com/golang/go/wiki/CodeReviewComments#pass-values) unless it is large (~20 fields or larger) or needs to be modified 
and returned by the callee.
- Don't ask for more than you need. In particular, if your function needs a user id, the
function should just ask for that id, and not the whole User object.

## Comments
Go is designed to be easy to read, but packages should always have comments and in most
cases, functions should also have comments. Guidelines :
- [General](https://blog.golang.org/godoc-documenting-go-code)
- [For package comments](https://github.com/golang/go/wiki/CodeReviewComments#package-comments)


## "Populate" pointer function arguments vs functions returning values

We prefer having methods that return results instead of ones which change parameters with pointers. For example:

    func AddNumbers(a, b int) int {...}

...is better than:

    func Add(a, b int, sum *int) {...}

But when you absolutely need to have such a function (which will update a pointer value), call it `Populate...(target, ...)`:

    function PopulateQuote(quote *Quote, ...args...) {...}

## Controller methods

Controller methods returning html content should be names as:

    func (s *Controllers) <HTTP_METHOD>Description(e echo.Context) {...}

For example:

    func (s *Controllers) GETBookings(e echo.Context) {...}

Handlers which are used for Widgets (or API endpoints) which return JSON, must have an `A` before the http method, for example:

    func (s *Controllers) APOSTBooking(e echo.Context) {...}

Those handlers are expected to return jsons. The clients are expected to send the `Accept: application/json` header.
But, they should also end with `.json`. For example `/merchant/company.json`.

## Error handling

We use merry (<https://github.com/ansel1/merry>) errors. Merry is a wrapper around errors with some additional data:

* Stacktrace
* User message
* Messages(s)

When an error happen, this is the common way to handle it:

    if err != nil {
        return merry.Wrap(err)
    }

If there are any other informations which will help the developer figure out the problem that occured, add messages:

    return merry.Wrap(err).Appendf("sql=%s", sql).Appendf("params=%#v", params)

The same error can have multiple errors (they are basically just appended with a colon). These errors are for
the programmer and will be logged in an echo middleware.

There are also **user errors**. They will be returned to the client, and they **must not** contain any debugging
informations (variables, stacktrace). User errors should be set only in controllers or in helper around controllers!

    if err != nil {
        return merry.Wrap(err).WithUserMessage(i18n.ERROR_LOGIN).WithHTTPCode(http.StatusUnauthorized)
    }

The rules for setting **HTTP code** are similar to user messages: It should be set only in controllers. If no http error
is set, then the response will be a `500 Internal Server Error`. Make sure that you have proper http responses for
most errors!

## Http code

Common uses:

* Validation errors: 400
* Entity created: 201
* Empty response: 204

## Database transactions

When a function has sql operations that must be run in a transaction **don't** start a transaction in the middle of the
function and commit/rollback there. For example:

    func (bs *BookingService) DeleteBooking(b *model.Booking) (err error) {
        err = bs.TransactionExecutor.ExecuteInTransaction(func(tx *sqlx.Tx) error {
            err = bs.DeleteBookingTx(tx, b)
            return
        })
        return
    }

    func (bs *BookingService) DeleteBookingTx(tx *sqlx.Tx, b *model.Booking) error {
        ...
    }

Note that the `DoSomething(...)` function must have named return values (in the example `booking` and `err`), otherwise
the deferred function cannot change/set the return values.
