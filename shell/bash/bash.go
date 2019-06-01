package bash

import(
  "runtime"
  "os/exec"
  "context"
  "time"
  "fmt"
  "github.com/pkg/errors"
  "log"
  "os"
)

type Bash struct {
  path string
  timeout time.Duration
}

func init() {
  log.SetOutput(os.Stderr)
}

func New(timeout time.Duration) *Bash {

  return &Bash{
            path: "/bin/bash",
            timeout: timeout,
        }
}

func (s *Bash) Execute(script []byte, args ...string) ([]byte, error){
  // start of call
  start := time.Now()

  // caller-id
  pc, _, _ , ok := runtime.Caller(1)
  details := runtime.FuncForPC(pc)
  caller := ""
  if ok && details != nil{
    caller = fmt.Sprintf("%s", details.Name())
  }

  // end of clal
  defer func(start time.Time, caller string) {
    // end of call
    end := time.Now()
    log.Printf(">> BASH >> exec caller %q took %d us\n", caller, (end.Sub(start).Nanoseconds()/ 1000))
  }(start, caller)


  ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
  defer cancel()

  cmd := exec.CommandContext(ctx, s.path, append([]string{"-c", string(script)}, args...)...)
  out, err := cmd.Output()
  if err != nil {
    if ctx.Err() != nil {
      return nil, errors.Wrap(err,fmt.Sprint(ctx.Err()))
    }
    return nil, err
  }

  return out, nil

}
