package cmd

import (
	"strings"

	"github.com/SAP/jenkins-library/pkg/command"
	"github.com/SAP/jenkins-library/pkg/log"
	"github.com/SAP/jenkins-library/pkg/telemetry"
	"github.com/pkg/errors"
)

func uiVeri5ExecuteTests(config uiVeri5ExecuteTestsOptions, telemetryData *telemetry.CustomData) {
	// for command execution use Command
	c := command.Command{}
	// reroute command output to logging framework
	c.Stdout(log.Writer())
	c.Stderr(log.Writer())

	// error situations should stop execution through log.Entry().Fatal() call which leads to an os.Exit(1) in the end
	err := runUIVeri5(&config, &c)
	if err != nil {
		log.Entry().WithError(err).Fatal("step execution failed")
	}
}

func runUIVeri5(config *uiVeri5ExecuteTestsOptions, command command.ExecRunner) error {
	envs := []string{"NPM_CONFIG_PREFIX=" + config.NpmConfigPrefix}
	envs = append(envs, "TARGET_SERVER_URL="+config.TestServerURL)
	command.SetEnv(envs)

	installCommandTokens := strings.Split(config.InstallCommand, " ")
	err := command.RunExecutable(installCommandTokens[0], installCommandTokens[1:]...)
	if err != nil {
		log.Entry().WithError(err).WithField("command", config.InstallCommand).Fatal("failed to execute install command")
		return err
	}

	if config.TestOptions != "" {
		return errors.Errorf("parameter testOptions no longer supported, please use runOptions parameter instead.")
	}
	err = command.RunExecutable(config.RunCommand, config.RunOptions...)
	if err != nil {
		log.Entry().WithError(err).WithField("command", config.RunCommand).Fatal("failed to execute run command")
		return err
	}
	return nil
}