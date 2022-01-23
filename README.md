# A lightweight logging interface for Go

## Feature highlights
* Minimal, non-opinionated interface
* **B**ring **Y**our **O**wn **B**ackend

## Usage

Add the module as a dependency to your project
```
go get github.com/siliconbrain/gologlite
```

As `gologlite` contains only a minimal API for logging, it's recommended to *create your own custom `log` package* in your project and:
* use alises to re-publish necessary types and functions of `gologlite` and any other modules/packages extending functionality herein
* implement your own helpers and utilities to best adapt to the needs of your project

The rest of your project should import this custom `log` package.
