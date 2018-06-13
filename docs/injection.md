Taken from Globio. modify.

# Dependency injection

For dependency injection we use <https://github.com/facebookgo/inject> with a small wrapper around it to make objects initializable.

## Lifecycle of a Dependency injection object

Objects which are part of dependency injection are singletons and contain important logic code for the application.
Objects which are instantiated, used and then left to the GC to be cleaned are not part of the depenedency injection.

## Example

If we have a UserControllers object (containing user echo handler methods), and this object need to use the logger and emailer, then it will be declared as:

    type UserControllers struct {
        Logger       log.Logger          `inject:""`
        EmailService *email.EmailService `inject:""`
    }
    
All fields used in DI must point be pointers. The reason why `log.Logger` dont have an `*` is that it's a interface (interfaces are pointers by default).

If `UserControllers` contains an embedded struct which, again contains fields to be injected, then we need:

    type UserControllers struct {
        BaseControllers `inject:"inline"`
        
        Logger       log.Logger          `inject:""`
        EmailService *email.EmailService `inject:""`
    }

### Declaration

Objects must be declared in `di.go`:

	di := diwrapper.New().
		WithObject(store).
		WithObject(cfg).
		WithObject(logger).
		InitializeGraph()
    defer di.Stopper()

When `InitializeGraph()` is called then all the DI fields are injected.

### Object initialization

If the object has some custom initialization (for example caching something from the database or setting up external API keys) code it must have the following method:

    func (o *Object) Init() error {...}
    
For example, a common usage is to setup an API key:

    type ExternalService struct {
        Config config.Config `inject:"""`
        
        apiKey string
    }
    
    func (s *ExternalService) Init() error {
        s.apiKey = s.Config.GetString("external.service.api.key")
        if len(s.apiKey) == 0 {
            return errors.New("No API key found for external API..")
        }
        return nil
    }

All `Init()` methods are called on startup **after** all the fields marked with `inject:"..."` are filled. 
If any of the `Init()` methods returns an error, the app will panic and report the exact objects where initialization failed.

### Named dependencies

If we want multiple instances of one object, they must be named:

	di := diwrapper.New().
		WithNamedObject("db_logger", new(NullLogger)).
		WithNamedObject("payment_logger", new(PaymentLogger)).
		WithNamedObject("controller_logger", new(StackdriverLogger)).

Objects that will use logging can request one (or more) specific loggers:

    type PayoneerService struct {
        Logger log.Logger `inject:"payment_logger""`
    }

or:

    type UserControllers struct {
        Logger log.Logger `inject:"controller_logger""`
    }

or (in the case where an object needs multiple loggers):

    type DoSomethingService struct {
        DbLogger log.Logger `inject:"db_logger""`
        Logger log.Logger `inject:"controller_logger""`
    }
