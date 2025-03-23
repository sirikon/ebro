## Breaking Changes

### New label resolution logic

Previously, child tasks would inherit parent tasks labels already resolved with the parent's environment. Now, child tasks inherit the parent's labels unresolved and are resolved using the child task's environment.

```yaml
tasks:
  parent:
    abstract: true
    environment:
      FOOD: pizza
    labels:
      food: $FOOD
  child:
    extends: [parent]
    environment:
      FOOD: paella
```

```yaml
# ebro -i

# before
:child:
  environment:
    FOOD: paella
  labels:
    food: pizza

# after
:child:
  environment:
    FOOD: paella
  labels:
    food: paella
```

## Miscellanea

Updated Go dependencies.
