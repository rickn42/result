# result
simple promise type structure

create
```go
r := NewResult()
```

set
```go 
r.Set(result-value)
```

wait
```go
select {
case <-r.Done():
    // result done. do something.
default:
    // result is not done.
}
```

result value 
```go
// this wait until done. 
val, err := r.Value() 
// result.Set(error) => (nil, error)
// result.Set(not-error-value) => (that-value, nil)
```

