crudo
----------------

A simple CRUD (Create, Read, Update, Delete) library for simple management of data.

## Features

- Based on the repository pattern, it abstracts the data access layer and provides a simple interface for managing data.
- Criteria (specification) pattern for querying data without bloating the repository implementation.
- MongoDB and InMemory (unsafe slice) implementation for the repository pattern.

## WIP

- Id value object is opinionated to provide simple operation of the library, it provides creation/updation hooks
  and get/set for the id. It also has a WIP of compound ids.
- InMemory with safe slice implementation.
- More tests coverage and examples.

## Installation

```bash
go get github.com/davfer/crudo
```

## Usage

model.go:

```go
package main

import (
	"github.com/davfer/crudo"
	"github.com/davfer/crudo/entity"
)

type User struct {
	entity.Id
	Slug string
	Name string
	Age  int
}

func (t *User) GetId() entity.Id {
	return t.Id
}

func (t *User) SetId(id entity.Id) error {
	t.Id = id
	return nil
}

func (t *User) GetResourceId() (string, error) {
	return t.Slug, nil
}

func (t User) SetResourceId(s string) error {
	t.Slug = s
	return nil
}

func (t User) PreCreate() error {
	return nil
}

func (t User) PreUpdate() error {
	return nil
}

```

main.go:

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/davfer/crudo"
)

func main() {
	
}
```


