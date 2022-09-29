

[Organize your code with Go packages](https://blog.learngoprogramming.com/code-organization-tips-with-packages-d30de0d11f46)

*  Smaller programs may not need many packages
*  Keep your packages small
*  Put related packages into sub-directories
*  Put tests into the same directory, Put the data needed for testing into testdata directories as a sub-directory
```
miner/
    miner.go
    miner_test.go
    testdata/
        hashes.data
        .
```


