package depot_tools_auth

import (
	"context"

	"go.skia.org/infra/go/exec"
	"go.skia.org/infra/go/git"
	"go.skia.org/infra/go/skerr"
	"go.skia.org/infra/go/sklog"
)

// GitConfig sets up the git config for depot tools authentication.
//
// Note: This requires `git-credential-luci` to be present and in PATH.
func GitConfig(ctx context.Context) error {
	sklog.Infof("Setting git configuration for depot tools auth...")
	gitExec, err := git.Executable(ctx)
	if err != nil {
		return skerr.Wrap(err)
	}
	if _, err := exec.RunCwd(ctx, ".", gitExec, "config", "--global", "credential.https://*.googlesource.com.helper", "luci"); err != nil {
		return skerr.Wrap(err)
	}
	if _, err := exec.RunCwd(ctx, ".", gitExec, "config", "--global", "--unset", "http.cookiefile"); err != nil {
		// Ignore the error; "git config --unset" will exit with a non-zero code
		// if the setting doesn't exist and therefore couldn't be removed.
		sklog.Warning("'git config --unset http.cookiefile' exited with non-zero code. Ignoring.")
	}
	// Read back gitconfig.
	out, err := exec.RunCwd(ctx, ".", gitExec, "config", "--list", "--show-origin")
	if err != nil {
		return skerr.Wrapf(err, "Failed to read git config: %s", out)
	}
	sklog.Infof("Created git configuration:\n%s", out)
	return nil
}
