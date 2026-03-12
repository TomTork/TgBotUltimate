package neuro

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"syscall"
	"time"
)

func Init(ctx context.Context) error {
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		return fmt.Errorf("runtime.Caller failed")
	}

	projectRoot := filepath.Clean(filepath.Join(filepath.Dir(thisFile), "..", ".."))

	if _, err := os.Stat(filepath.Join(projectRoot, "training", "api.py")); err != nil {
		return fmt.Errorf("cannot find training/api.py from dir=%s: %w", projectRoot, err)
	}
	cmd := exec.CommandContext(
		ctx,
		"python3",
		"training/api.py",
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	if err := cmd.Start(); err != nil {
		return err
	}

	pid := cmd.Process.Pid
	log.Printf("[neuro] python started pid=%d dir=%s\n", pid, projectRoot)

	go func() {
		<-ctx.Done()
		_ = syscall.Kill(-pid, syscall.SIGTERM)
		time.Sleep(2 * time.Second)
		_ = syscall.Kill(-pid, syscall.SIGKILL)
	}()

	err := cmd.Wait()

	if ctx.Err() != nil {
		return nil
	}

	if err != nil {
		return fmt.Errorf("[neuro] python exited: %w", err)
	}

	return errors.New("[neuro] python exited unexpectedly")
}
