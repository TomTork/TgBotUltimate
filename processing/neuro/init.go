package neuro

import (
	"log"
	"os/exec"
	"syscall"
)

func InitPython() error {
	cmd := exec.Command("python3", "main.py")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}
	if err := cmd.Start(); err != nil {
		return err
	}
	log.Printf("Python запущен с PID %d", cmd.Process.Pid)
	return nil
}
