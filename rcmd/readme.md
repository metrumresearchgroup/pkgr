# Readme

There is an additional integration test for testing parsing R outputs, in which R 3.5.2 is required on the system

To run it, use the `R` tag

```
go test -tags=R
```

This should eventually be refactored in some way, potentially using env variables or some other mechanism to 
set the R versions on the system to determine whether to run the test and/or what the output should be.