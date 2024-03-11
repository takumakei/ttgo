package ttgo

import (
	"bufio"
	"bytes"
	_ "embed"
	"encoding/json"
	"io"
	"log/slog"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/takumakei/ttgo/headline"
)

//go:embed help.txt
var helpTxt string

var Command = &cobra.Command{
	Use:   "ttgo [flags] <FILE>...",
	Short: headline.Get(helpTxt),
	Long:  helpTxt,
	RunE:  Run,
}

var (
	FlagTmpl string
	FlagData []string
)

func init() {
	flags := Command.Flags()
	flags.SortFlags = false
	flags.StringVarP(&FlagTmpl, "tmpl", "t", FlagTmpl, "template `file` (required)")
	flags.StringArrayVarP(&FlagData, "data", "d", FlagData, "input `data`")
}

func Run(cmd *cobra.Command, args []string) error {
	if FlagTmpl == "" {
		return pflag.ErrHelp
	}

	tmpl, err := NewTemplate(FlagTmpl)
	if err != nil {
		return err
	}

	if len(FlagData) == 0 && len(args) == 0 {
		return Exec(tmpl, bufio.NewReader(os.Stdin), Input{Type: STDIN})
	}

	for _, e := range FlagData {
		if err := Exec(tmpl, bytes.NewBuffer([]byte(e)), Input{Type: DATA}); err != nil {
			return err
		}
	}

	for _, e := range args {
		if err := ExecFile(tmpl, e, Input{Type: FILE, File: e}); err != nil {
			return err
		}
	}

	return nil
}

func Exec(tmpl *Template, r io.Reader, it Input) error {
	data, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	input := &Context{Input: it, Data: string(data)}
	slog.Debug("Exec", slog.Any("input", input))
	result, err := tmpl.Execute(input)
	if err != nil {
		return err
	}
	_, err = os.Stdout.Write([]byte(result))
	return err
}

func ExecFile(tmpl *Template, file string, it Input) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()
	return Exec(tmpl, bufio.NewReader(f), it)
}

type Context struct {
	Input Input  `json:"input"`
	Data  string `json:"data"`
}

func (c *Context) String() string {
	j, _ := json.Marshal(c)
	return string(j)
}

type Input struct {
	Type InputType `json:"type"`
	File string    `json:"file"`
}

type InputType string

const (
	STDIN InputType = "stdin"
	DATA  InputType = "data"
	FILE  InputType = "file"
)
