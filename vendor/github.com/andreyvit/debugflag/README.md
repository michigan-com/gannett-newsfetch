# debugflag

Golang library to enable optional behavior at runtime via `DEBUG` environment variable. Use it for:

* very verbose logging
* time-consuming checks
* relaxing limits and bypassing restrictions
* other behaviors undesirable in production apps and even during general developement.

Similar in spirit to [npm debug package](https://www.npmjs.com/package/debug), although this one is agnostic to logging; you can use it to enable any conditional behavior at runtime.


## Usage

Set `DEBUG` environment variable to a comma-separated list of debug flags you want to enable:

    DEBUG=json,net ./foo

Use `IsEnabled(flag string)` to check if a given flag is enabled:

    import (
        "github.com/andreyvit/debugflag"
    )

    func process() {
        json := LoadHugeJSON()
        if debugflag.IsEnabled("json") {
            // dump the entire json to console here
        }
    }


### Namespaced flags

Optionally, use colons to namespace your flags, and enable them in bulk:

    DEBUG=foo:net ./foo

and:

    func GetArticles() {
        ...
        if debugflag.IsEnabled("foo:net:articles") {
            // dump articles json
        }
        ...
    }

    func GetComments() {
        ...
        if debugflag.IsEnabled("foo:net:comments") {
            // dump articles json
        }
        ...
    }


### Disabling flags

You can disable certain flags by prefixing them with `-` or `!`. For example, to dump all network requests except for articles, you can use:

    DEBUG=foo:net,-foo:net:comments ./foo


### Enabling all flags

Use `DEBUG=all` to enable all flags. You can make exceptions for specific flags:

    DEBUG=all,-foo:net:comments


### Specificity and precedence

The order of the flags in the `DEBUG` variable does _not_ matter; _specificity_ does matter, however, and more specific flags take precendence over less specific. For example, all of these:

    DEBUG=all,-foo,foo:net,-foo:net:comments
    DEBUG=-foo,-foo:net:comments,foo:net,all
    DEBUG=foo:net,all,-foo,-foo:net:comments

result in the same configuration:

* disabled: `foo`, `foo:whatever`, `foo:net:comments` and `foo:net:comments:whatever`;
* enabled: `foo:net`, its descendants (except for `foo:net:comments`) and all other unspecified flags.


## License

Copyright 2016, Andrey Tarantsov. Distributed under the [MIT license](LICENSE).
