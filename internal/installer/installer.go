package installer

import (
	"github.com/charmbracelet/log"
	"github.com/lunar-debian/wsl-dev-setup/internal/config"
	"github.com/lunar-debian/wsl-dev-setup/internal/utils"
)

type Context struct {
	Env     *utils.Environment
	Config  config.InstallConfig
	Logger  *log.Logger
	Backup  *utils.BackupManager
	DryRun  bool
}

type Step interface {
	Name() string
	Check(ctx *Context) (bool, error)
	Run(ctx *Context) error
	Rollback(ctx *Context) error
}

type Installer struct {
	logger    *log.Logger
	env       *utils.Environment
	config    config.InstallConfig
	steps     []Step
	backupMgr *utils.BackupManager
}

func New(logger *log.Logger, env *utils.Environment, cfg config.InstallConfig) *Installer {
	backupMgr := utils.NewBackupManager(logger)
	inst := &Installer{
		logger:    logger,
		env:       env,
		config:    cfg,
		backupMgr: backupMgr,
	}
	inst.registerSteps()
	return inst
}

func (i *Installer) registerSteps() {
	i.steps = []Step{
		NewSystemBaseStep(),
		NewFontsStep(),
		NewOMZStep(),
		NewP10kStep(),
		NewPluginsStep(),
		NewCLIToolsStep(),
		NewShellConfigStep(),
		NewOpencodeStep(),
		NewFinalizeStep(),
	}
}

func (i *Installer) Run() error {
	ctx := &Context{
		Env:     i.env,
		Config:  i.config,
		Logger:  i.logger,
		Backup:  i.backupMgr,
		DryRun:  i.config.DryRun,
	}

	i.logger.Info("Starting installation", "style", i.config.Style, "dry_run", i.config.DryRun)

	for _, step := range i.steps {
		i.logger.Info("Checking step", "step", step.Name())

		done, err := step.Check(ctx)
		if err != nil {
			i.logger.Error("Check failed", "step", step.Name(), "error", err)
			return err
		}

		if done {
			i.logger.Info("Step already done, skipping", "step", step.Name())
			continue
		}

		i.logger.Info("Running step", "step", step.Name())
		if err := step.Run(ctx); err != nil {
			i.logger.Error("Step failed, rolling back", "step", step.Name(), "error", err)
			i.rollbackSteps(ctx)
			return err
		}
		i.logger.Info("Step completed", "step", step.Name())
	}

	i.logger.Info("Installation completed successfully!")
	i.printNextSteps()
	return nil
}

func (i *Installer) rollbackSteps(ctx *Context) {
	i.logger.Info("Rolling back...")
	for j := len(i.steps) - 1; j >= 0; j-- {
		step := i.steps[j]
		if err := step.Rollback(ctx); err != nil {
			i.logger.Warn("Rollback failed", "step", step.Name(), "error", err)
		}
	}
	i.backupMgr.RestoreAll()
}

func (i *Installer) Uninstall() error {
	ctx := &Context{
		Env:     i.env,
		Config:  i.config,
		Logger:  i.logger,
		Backup:  i.backupMgr,
		DryRun:  i.config.DryRun,
	}

	i.logger.Info("Starting uninstall...")

	for j := len(i.steps) - 1; j >= 0; j-- {
		step := i.steps[j]
		if err := step.Rollback(ctx); err != nil {
			i.logger.Warn("Uninstall step failed", "step", step.Name(), "error", err)
		}
	}

	i.backupMgr.RestoreAll()
	i.logger.Info("Uninstall completed")
	return nil
}

func (i *Installer) printNextSteps() {
	i.logger.Info("Next steps:")
	i.logger.Info("  1. Restart your terminal or run: exec zsh")
	i.logger.Info("  2. Run 'p10k configure' to customize Powerlevel10k")
	i.logger.Info("  3. Open opencode to see the configured theme")
}