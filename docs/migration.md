# V6 Migration

This documents describes the migration logic and defines the rules that needs to be applied.

## Migration Logic

### Mutating objects

Add clone to mutating objects. To achieve this, detect cases like:

#### Mutation case 1

```golang
foo.Id = foo.NewId()
```

Replace with following:

```golang
cFoo := foo.Apply(&FooPatch{
    Id : model.NewId(),
})
```

#### Mutation case 2

```golang
bars := foo.Bars
bars[0] = newBar{}
```

Replace with following:

```golang
bars := foo.Bars() // will return a copy of the slice/map
bars[0] = newBar{}
```

### Field Access

We should replace every field access with method access.

```golang
id := channel.Id
```

```golang
id := channel.ID() // ensure value is a copy
```

### Object Creation

Replace all creation with initialization/patching.

```golang
foo := &model.Foo{
    Id: model.NewId(),
}
```

Replace with following:

```golang
foo := model.NewFoo(&FooInitializer{
    ID: model.NewID(), // ensure pointer values are cloned/deep-copied
})
```
