package converter

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type vpype struct {
	cmd    *exec.Cmd
	stderr io.ReadCloser
}

// defaultVpypeCommand is the default command to use for starting vpype from PATH
const defaultVpypeCommand = "vpype"

// VpypeOption is an option for the vpype converter
type VpypeOption func(*vpype)

// VpypeCommand sets the command to use for vpype
func VpypeCommand(cmd string) VpypeOption {
	return func(v *vpype) {
		v.cmd = exec.Command(cmd)
	}
}

// Vpype returns a new converter that uses vpype
func Vpype(opts ...VpypeOption) *vpype {
	v := &vpype{cmd: exec.Command(defaultVpypeCommand)}

	for _, opt := range opts {
		opt(v)
	}

	return v
}

// vpypeWriter is a writer that converts svg to hpgl
type vpypeWriter struct {
	svg    io.Reader
	cmd    *exec.Cmd
	config converterConfig
}

// vpypeWriter implements io.WriterTo
var _ io.WriterTo = &vpypeWriter{}

// Convert returns a writer that converts the svg to hpgl
func (v *vpype) Converter(svg io.Reader, opts ...Option) io.WriterTo {
	return &vpypeWriter{
		svg:    svg,
		cmd:    v.cmd,
		config: config(opts),
	}
}

// WriteTo converts the svg to hpgl and writes it to out
func (w *vpypeWriter) WriteTo(out io.Writer) (int64, error) {
	tmp, err := os.CreateTemp("", "plotq-*.svg")
	if err != nil {
		return 0, fmt.Errorf("could not create temporary file: %w", err)
	}
	defer os.Remove(tmp.Name())

	_, err = io.Copy(tmp, w.svg)
	if err != nil {
		return 0, fmt.Errorf("could not copy svg to temporary file: %w", err)
	}
	tmp.Close()

	w.cmd.Args = append(w.cmd.Args, commandArgs(tmp.Name(), w.config)...)

	stdoutPipe, err := w.cmd.StdoutPipe()
	if err != nil {
		return 0, fmt.Errorf("could not get stdout pipe: %w", err)
	}

	stderrPipe, err := w.cmd.StderrPipe()
	if err != nil {
		return 0, fmt.Errorf("could not get stderr pipe: %w", err)
	}

	if err := w.cmd.Start(); err != nil {
		return 0, fmt.Errorf("could not start command: %w", err)
	}

	written, err := io.Copy(out, stdoutPipe)
	if err != nil {
		return written, fmt.Errorf("could not copy from stdout to out: %w", err)
	}

	stderr, err := io.ReadAll(stderrPipe)
	if err != nil {
		return written, fmt.Errorf("could not read all from stderr: %w", err)
	}

	if w.cmd.Wait() != nil {
		// vpype prints out long tracebacks on stderr and only the last line is the actual error
		s := bufio.NewScanner(bytes.NewReader(stderr))
		for s.Scan() {
			err = errors.New(s.Text())
		}
		return written, fmt.Errorf("vpype %w", err)
	}

	return written, nil
}

// commandArgs returns the arguments for the vpype command
func commandArgs(svg string, cfg converterConfig) []string {
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

	return append(args, "--format", "hpgl", "-")
}
