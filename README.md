# wx
`wx`(.go) is a command-line application that retrieves [NOAA](http://www.weather.gov) weather data.

## Installation
Ensure Go is [installed](http://golang.org/doc/install.html) and your `PATH` includes the `$GOPATH/bin` directory

```
export PATH=$PATH:$GOPATH/bin
```

then install `wx` with

```
$ go get github.com/roualdes/wx
```

## Getting Started
Only two commands are available `forecast` and `current`.  Each command takes an optional argument `ZIPCODE`, a five digit ZIP code corresponding to the location of interest.  Try

```
$ wx current 40502
```

If you don't specify a ZIP code, `wx` uses 95926.  Print help with

```
$ wx help
```



