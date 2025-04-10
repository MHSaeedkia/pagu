package main

import (
	"os"
	"os/signal"
	"syscall"

	pagucmd "github.com/pagu-project/pagu/cmd"
	"github.com/pagu-project/pagu/config"
	"github.com/pagu-project/pagu/internal/engine"
	"github.com/pagu-project/pagu/internal/platforms/discord"
	"github.com/pagu-project/pagu/pkg/log"
	"github.com/spf13/cobra"
)

func runCommand(parentCmd *cobra.Command) {
	run := &cobra.Command{
		Use:   "run",
		Short: "Runs an instance of Pagu",
	}

	parentCmd.AddCommand(run)

	run.Run = func(cmd *cobra.Command, _ []string) {
		// load configuration.
		configs, err := config.Load(configPath)
		pagucmd.ExitOnError(cmd, err)

		// Initialize global logger.
		log.InitGlobalLogger(configs.Logger)

		// starting eng.
		eng, err := engine.NewBotEngine(configs)
		pagucmd.ExitOnError(cmd, err)

		eng.Start()

		bot, err := discord.NewDiscordBot(eng, configs.Discord, configs.BotName)
		pagucmd.ExitOnError(cmd, err)

		err = bot.Start()
		pagucmd.ExitOnError(cmd, err)

		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
		<-sigChan

		if err := bot.Stop(); err != nil {
			pagucmd.ExitOnError(cmd, err)
		}

		eng.Stop()
	}
}
