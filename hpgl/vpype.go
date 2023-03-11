package hpgl

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type vpype struct {
	cmd string
}

const defaultVpypeCommand = "vpype"

type VpypeOption func(*vpype)

func VpypeCommand(cmd string) VpypeOption {
	return func(v *vpype) {
		v.cmd = cmd
	}
}

func VpypeConverter(opts ...VpypeOption) *vpype {
	v := &vpype{cmd: defaultVpypeCommand}

	for _, opt := range opts {
		opt(v)
	}

	return v
}

func (v *vpype) Convert(svg io.Reader, opts ...ConvertOption) (hpgl io.ReadCloser, err error) {
	cfg := config(opts)

	tmp, err := os.CreateTemp("", "plotq-*.svg")
	if err != nil {
		return nil, fmt.Errorf("could not create temporary file: %w", err)
	}
	defer os.Remove(tmp.Name())

	_, err = io.Copy(tmp, svg)
	if err != nil {
		return nil, fmt.Errorf("could not copy svg to temporary file: %w", err)
	}
	tmp.Close()

	cmd := v.command(tmp.Name(), cfg)

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("could not get stdout pipe: %w", err)
	}

	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return nil, fmt.Errorf("could not get stderr pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("could not start command: %w", err)
	}

	stdout, err := io.ReadAll(stdoutPipe)
	if err != nil {
		return nil, fmt.Errorf("could not read all from stdout: %w", err)
	}

	stderr, err := io.ReadAll(stderrPipe)
	if err != nil {
		return nil, fmt.Errorf("could not read all from stderr: %w", err)
	}

	if cmd.Wait() != nil {
		s := bufio.NewScanner(bytes.NewReader(stderr))
		for s.Scan() {
			err = errors.New(s.Text())
		}
		return nil, fmt.Errorf("vpype %w", err)
	}

	return ioutil.NopCloser(bytes.NewReader(stdout)), nil
}

func (v *vpype) command(svg string, cfg converterConfig) *exec.Cmd {
	args := []string{"read", svg, "write"}

	if cfg.landscape {
		args = append(args, "--landscape")
	}

	if cfg.pagesize != "" {
		args = append(args, "--page-size", strings.ToLower(cfg.pagesize))
	}

	if cfg.device != "" {
		args = append(args, "--device", strings.ToLower(cfg.device))
	}

	if cfg.velocity != 0 {
		args = append(args, "--velocity", strconv.Itoa(int(cfg.velocity)))
	}

	args = append(args, "--format", "hpgl", "-")

	return exec.Command(v.cmd, args...)
}
