# Doppler secrets subsitution

My First *go* at coding **Go** 

A little project with [Doppler](https://www.doppler.com/)

## What it does

Very simple api/rest function to retrieve secrets stored in Doppler and
replaces variable expressions (e.g. `${DATABASE_URL}` ) in static files with the
respective Doppler secret.

Variable Expression Formats Supported:

 * dollar-curly      e.g. `${MYVAR}` - Default

I need to fix:
 * dollar            e.g. `$MYVAR`
 * handlebars        e.g. `{{MYVAR}}`
 * dollar-handlebars e.g. `${{MYVAR}}`

## Prerequesites
Here's what you will need to have

 * [Golang](https://go.dev/dl/) > v1.18
 * [Doppler Service Token](https://docs.doppler.com/docs/service-tokens)

## Installation

*WIP*

## Usage

`DOPPLER_TOKEN` has to be initialised.

```bash
DOPPLER_TOKEN="dp.st.dev_xxxxxxxxxxxxxxxxx" go run doppler-sub.go ./files ./export dollar-curly
DOPPLER_TOKEN="dp.st.dev_xxxxxxxxxxxxxxxxx" go run doppler-sub.go ./files ./export
```

Example output
```bash
Variable expression format to target:  dollar-curly i.e. \${([A-Z_]{1,}[A-Z0-9_].*?)}
Reading files/myfile
        Secrets Matched:
                 I_DONT_EXIST           ✗
                 I_DONT_EXIST           ✗
                 LOGGING                ✔
                 STRIPE_KEY             ✔
                 PRIVATE_KEY            ✔
                 FEATURE_FLAGS          ✔
                 DOPPLER_CONFIG         ✔
                 DOPPLER_ENVIRONMENT    ✔
                 DOPPLER_PROJECT        ✔
                 DOPPLER_ENVIRONMENT    ✔
                 DOPPLER_ENVIRONMENT    ✔
                 DOPPLER_ENVIRONMENT    ✔
                 DOPPLER_ENVIRONMENT    ✔
        Total variables matched: 11
        Secrets written to ./export/myfile
Reading files/orginal.txt
        Secrets Matched:
                 DB_URL                 ✔
                 DOPPLER_ENVIRONMENT    ✔
                 DOPPLER_ENVIRONMENT    ✔
                 LOGGING                ✔
                 STRIPE_KEY             ✔
                 PRIVATE_KEY            ✔
                 FEATURE_FLAGS          ✔
                 DOPPLER_CONFIG         ✔
                 DOPPLER_ENVIRONMENT    ✔
                 DOPPLER_PROJECT        ✔
                 I_DONT_EXIST           ✗
        Total variables matched: 10
        Secrets written to ./export/orginal.txt
```

## Todo

 * **Get better at writing GOlang**
 * Get the following Variable expressions working:
   * dollar            e.g. `$MYVAR`
   * handlebars        e.g. `{{MYVAR}}`
   * dollar-handlebars e.g. `${{MYVAR}}`

* [Command-Line Flags](https://gobyexample.com/command-line-flags)
* Use [tabwriter](https://pkg.go.dev/text/tabwriter) for prettier output
* Optimize for speed for large batch processing