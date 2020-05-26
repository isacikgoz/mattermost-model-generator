# V6 Migration

This documents describes the migration logic and defines the rules that needs to be applied.

## Migration Logic

### Mutating objects

Add clone to mutating objects. Detect:
1. Direct assignments, i.e. `obj.Id = model.NewId()`
2. Inderect assignments, i.e. `obj.Stuff[0] = newStuff` and `obj.Stuff["test"] = "test2"`

If any of the above cases are detected, rewrite the function so that the model parameter passed to it (for example `channel *model.Channel`) is renamed (for example `_channel *model.Channel`) and add a clone to the first line: `channel := _channel.Clone()`
This will make sure that all following assignments will be applied to the clone of the model and not the origin.

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
