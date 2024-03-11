package logger

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/gobwas/glob"
)

var root *slog.Logger = initRoot()

func init() {
	slog.SetDefault(root)
}

func Root() *slog.Logger {
	return root
}

func initRoot() *slog.Logger {
	var (
		addSource = true
		slogText  = false
		out       = io.Writer(os.Stderr)
		level     = slog.LevelInfo
		err       error
		levels    []*packageLevel
	)

	if s, ok := os.LookupEnv("SLOG_ADD_SOURCE"); ok {
		addSource, err = strconv.ParseBool(s)
		if err != nil {
			panic(fmt.Sprintf("SLOG_ADD_SOURCE=%q, cannot parse as bool", s))
		}
	}

	if s, ok := os.LookupEnv("SLOG_TEXT"); ok {
		slogText, err = strconv.ParseBool(s)
		if err != nil {
			panic(fmt.Sprintf("SLOG_TEXT=%q, cannot parse as bool", s))
		}
	}

	if s, ok := os.LookupEnv("SLOG_OUTPUT"); ok {
		switch s {
		case "out", "stdout":
			out = os.Stdout
		case "err", "stderr":
			out = os.Stderr
		default:
			panic(fmt.Sprintf("SLOG_OUTPUT=%q, must be stdout or stderr", s))
		}
	}

	if s, ok := os.LookupEnv("SLOG_LEVEL"); ok {
		re := regexp.MustCompile(`^([^=]+)=(.*)$`)
		for _, e := range strings.Split(s, ",") {
			e = strings.TrimSpace(e)
			if m := re.FindStringSubmatch(e); len(m) > 0 {
				lhs, rhs := m[1], m[2]
				g, err := glob.Compile(lhs)
				if err != nil {
					panic(fmt.Sprintf("SLOG_LEVEL=..,%s,..; %q; %v", e, lhs, err))
				}
				var level slog.Level
				if err := level.UnmarshalText([]byte(rhs)); err != nil {
					panic(fmt.Sprintf("SLOG_LEVEL=..,%s,..; %q; %v", e, rhs, err))
				}
				levels = append(levels, &packageLevel{glob: g, level: level})
			} else {
				if err := level.UnmarshalText([]byte(e)); err != nil {
					panic(fmt.Sprintf("SLOG_LEVEL=..,%s,..; %v", e, err))
				}
			}
		}
	}

	opts := &slog.HandlerOptions{
		AddSource: addSource,
		Level:     level,
	}
	var h slog.Handler
	if slogText {
		h = slog.NewTextHandler(out, opts)
	} else {
		h = slog.NewJSONHandler(out, opts)
	}
	return slog.New(&handler{
		handler: h,
		levels:  levels,
	})
}
