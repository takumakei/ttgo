package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/spf13/cobra"
	_ "github.com/takumakei/ttgo/logger"
	"github.com/takumakei/ttgo/ttgo"
)

func main() {
	ttgo.Command.Version = version
	Run(context.Background(), ttgo.Command)
}

func Run(ctx context.Context, cmd *cobra.Command) {
	if err := run(ctx, cmd); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, cmd *cobra.Command) (err error) {
	slog.Log(ctx, slog.LevelDebug-2, "START", slog.String("cmd", cmd.CommandPath()))
	defer func() {
		slog.Log(ctx, slog.LevelDebug-2, "EXIT", slog.String("cmd", cmd.CommandPath()), slog.Any("err", err))
	}()

	cobra.EnableCaseInsensitive = true
	cobra.EnableCommandSorting = false
	cobra.EnablePrefixMatching = true

	cmd.SilenceErrors = true
	cmd.SilenceUsage = true
	return cmd.ExecuteContext(ctx)
}
