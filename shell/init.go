package shell

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/LNKLEO/OMP/log"
	"github.com/LNKLEO/OMP/runtime"
	"github.com/LNKLEO/OMP/runtime/path"
	"github.com/google/uuid"
)

const (
	noExe = "echo \"Unable to find OMP executable\""
)

func getExecutablePath(env runtime.Environment) (string, error) {
	executable, err := os.Executable()
	if err != nil {
		return "", err
	}

	if env.Flags().Strict {
		return path.Base(executable), nil
	}

	// On Windows, it fails when the excutable is called in MSYS2 for example
	// which uses unix style paths to resolve the executable's location.
	// PowerShell knows how to resolve both, so we can swap this without any issue.
	if env.GOOS() == runtime.WINDOWS {
		executable = strings.ReplaceAll(executable, "\\", "/")
	}

	return executable, nil
}

func Init(env runtime.Environment, feats Features) string {
	shell := env.Flags().Shell

	switch shell {
	case PWSH, PWSH5:
		executable, err := getExecutablePath(env)
		if err != nil {
			return noExe
		}

		var additionalParams string
		if env.Flags().Strict {
			additionalParams += " --strict"
		}

		var command, config string

		switch shell {
		case PWSH, PWSH5:
			command = "(@(& %s init %s --config=%s --print%s) -join \"`n\") | Invoke-Expression"
		}

		config = quotePwshStr(env.Flags().Config)
		executable = quotePwshStr(executable)

		return fmt.Sprintf(command, executable, shell, config, additionalParams)
	case ZSH, BASH, CMD:
		return PrintInit(env, feats, nil)
	default:
		return fmt.Sprintf(`echo "%s is not supported by OMP"`, shell)
	}
}

func PrintInit(env runtime.Environment, features Features, startTime *time.Time) string {
	executable, err := getExecutablePath(env)
	if err != nil {
		return noExe
	}

	shell := env.Flags().Shell
	configFile := env.Flags().Config
	sessionID := uuid.NewString()

	var script string

	switch shell {
	case PWSH, PWSH5:
		executable = quotePwshStr(executable)
		configFile = quotePwshStr(configFile)
		sessionID = quotePwshStr(sessionID)
		script = pwshInit
	case ZSH:
		executable = QuotePosixStr(executable)
		configFile = QuotePosixStr(configFile)
		sessionID = QuotePosixStr(sessionID)
		script = zshInit
	case BASH:
		executable = QuotePosixStr(executable)
		configFile = QuotePosixStr(configFile)
		sessionID = QuotePosixStr(sessionID)
		script = bashInit
	case CMD:
		executable = escapeLuaStr(executable)
		configFile = escapeLuaStr(configFile)
		sessionID = escapeLuaStr(sessionID)
		script = cmdInit
	default:
		return fmt.Sprintf("echo \"No initialization script available for %s\"", shell)
	}

	init := strings.NewReplacer(
		"::OMP::", executable,
		"::CONFIG::", configFile,
		"::SHELL::", shell,
		"::SESSION_ID::", sessionID,
	).Replace(script)

	shellScript := features.Lines(shell).String(init)

	if !env.Flags().Debug {
		return shellScript
	}

	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("\n%s %s\n", log.Text("Init duration:").Green().Bold().Plain(), time.Since(*startTime)))

	builder.WriteString(log.Text("\nScript:\n\n").Green().Bold().Plain().String())
	builder.WriteString(shellScript)

	builder.WriteString(log.Text("\n\nLogs:\n\n").Green().Bold().Plain().String())
	builder.WriteString(env.Logs())

	return builder.String()
}
