

systemd resource limit
https://fredrikaverpil.github.io/2016/04/27/systemd-and-resource-limits/


go 打印limit:
https://golang.org/pkg/syscall/#Getrlimit
var limit syscall.Rlimit
syscall.Getrlimit(syscall.RLIMIT_NOFILE, &limit)
log.Infof("limit got %+v", limit)