# valdy

[![golangci-lint](https://github.com/k1gabyt0/valdy/actions/workflows/golangci-lint.yaml/badge.svg?branch=5-add-golangci-lint-to-project)](https://github.com/k1gabyt0/valdy/actions/workflows/golangci-lint.yaml)

## Description

valdy is a go package, that provides simple and flexible validations.

⚡️ The project is in rapid stage of development, issues and critics are welcome.

### Features

- reflection free
- simple minded
- generics usage
- multierrors(using [erry](https://github.com/k1gabyt0/erry))

## Table of contents

- [valdy](#valdy)
  - [Description](#description)
    - [Features](#features)
  - [Table of contents](#table-of-contents)
  - [Requirements](#requirements)
  - [Getting started](#getting-started)
    - [Installation](#installation)
    - [Usage](#usage)
      - [Simple validation](#simple-validation)
      - [Simple validation with custom error](#simple-validation-with-custom-error)
  - [Alternatives](#alternatives)

## Requirements

go 1.18+

## Getting started

### Installation

```bash
go get github.com/k1gabyt0/valdy
```

### Usage

#### Simple validation

```go
type person struct {
 name     string
 age      int
 children []person
}

var isAdult = func(p person) valdy.Checker[person] {
  return valdy.NewRule(
    fmt.Sprintf("%s must be adult, but his age is %d", p.name, p.age),
     func(p person) bool {
       return p.age >= 18
     },
  )
}

func main() {
  var validator valdy.Validator[person]

  john := person{
    name: "John",
    age:  17,
  }

  err := validator.Validate(john,
    isAdult,
  )
  fmt.Println(err)

  fmt.Println()

  if errors.Is(err, valdy.ErrValidation) {
    fmt.Println("This is general validation error")
  }
}
```

**Output**:

```text
validation has failed:
  John must be adult, but his age is 17

This is general validation error
```

#### Simple validation with custom error

```go
type person struct {
 name     string
 age      int
 children []person
}

var errIsAdult = errors.New("age verification error")
var isAdult = func(p person) valdy.Checker[person] {
  return valdy.NewRule(
    fmt.Sprintf("%s must be adult, but his age is %d", p.name, p.age),
     func(p person) bool {
       return p.age >= 18
     },
  ).WithError(errIsAdult)
}

func main() {
  var validator valdy.Validator[person]

  john := person{
    name: "John",
    age:  17,
  }

  err := validator.Validate(john,
    isAdult,
  )
  fmt.Println(err)

  fmt.Println()

  if errors.Is(err, valdy.ErrValidation) {
    fmt.Println("This is general validation error")
  }
  if errors.Is(err, errIsAdult) {
    fmt.Println("This is age verification error")
  }
}
```

**Output**:

```text
validation has failed:
  age verification error:
  John must be adult, but his age is 17

This is general validation error
This is age verification error
```

## Alternatives

- [ozzo-validation](https://github.com/go-ozzo/ozzo-validation)
- [govalidator](https://github.com/asaskevich/govalidator)
