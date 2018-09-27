package api

import (
	"context"
	"log"
	"os/exec"
	"time"
)

type Cmd struct {
	*exec.Cmd
	StartTime  time.Time
	FinishTime time.Time
}

func Command(arg string) *Cmd {
	log.Printf("registering command %s\n", arg)
	return &Cmd{Cmd: exec.Command("sh", "-c", arg), StartTime: time.Now()}
}

func CommandContext(ctx context.Context, arg string) *Cmd {
	log.Printf("registering command-context %s\n", arg)
	return &Cmd{Cmd: exec.CommandContext(ctx, "bash", "-c", arg), StartTime: time.Now()}
}

func Start(cmd *Cmd) *Cmd {
	log.Printf("starting command %s\n", cmd.Args)
	if err := cmd.Start(); err != nil {
		log.Printf("starting failed on %s\n", cmd.Args)
		log.Println(err)
		return cmd
	}
	return cmd
}

func Wait(cmd *Cmd) error {
	err := cmd.Wait()
	log.Println("done with command")
	cmd.FinishTime = time.Now()
	return err
}
