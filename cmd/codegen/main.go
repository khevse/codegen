package main

import (
	"log"

	"github.com/khevse/codegen/internal/command/interface_creator"
	"github.com/khevse/codegen/internal/command/object_test_wrapper"
	"github.com/khevse/codegen/internal/pkg/command"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := cobra.Command{
		Use:  "generate",
		Args: cobra.ArbitraryArgs,
	}

	for _, cmd := range []command.Command{
		interface_creator.New(),
		object_test_wrapper.New(),
	} {
		childCmd := &cobra.Command{
			Use:   cmd.Name(),
			Short: cmd.ShortName(),
			RunE: func(*cobra.Command, []string) error {
				return cmd.Execute()
			},
		}
		if err := cmd.InitFlags(childCmd); err != nil {
			log.Fatalf("init flags(command=%s): %s", cmd.Name(), err)
		}

		rootCmd.AddCommand(childCmd)
	}

	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("execute command: %s", err)
	}
}
