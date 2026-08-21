# BUG_REPRO

The following failures were observed while validating the initial project state.
Each section records what failed, how to reproduce it, and the complete command output.
They are preserved intentionally; only failing build gates are omitted from the generated Dockerfile.

## Failure 1: Go test (.)

- Observed problem: `Go test (.)` failed in the initial project state.
- Working directory: `.`
- Command: `cd /app && GOTOOLCHAIN=local GOPROXY=off GOSUMDB=off go test -count=1 ./...`
- Exit status: `1`

```text
--- FAIL: TestBusiness36Regression (0.01s)
    integration_test.go:30: archive should preserve independent confirmed state
FAIL
FAIL	gestureflame	0.015s
?   	gestureflame/clock	[no test files]
?   	gestureflame/cmd/flamectl	[no test files]
ok  	gestureflame/cli	0.007s
ok  	gestureflame/domain	0.003s
ok  	gestureflame/httpapi	0.013s
ok  	gestureflame/report	0.002s
ok  	gestureflame/service	0.019s
ok  	gestureflame/store	0.011s
FAIL
```

## Architecture reproduction

### linux/amd64
- Go toolchain version: exit `0`
- Go build (.): exit `0`
- Go test (.): exit `1`
- Go run smoke (cmd/flamectl): exit `0`
### linux/arm64
- Go toolchain version: exit `0`
- Go build (.): exit `0`
- Go test (.): exit `1`
- Go run smoke (cmd/flamectl): exit `0`
