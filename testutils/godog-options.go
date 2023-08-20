package testutils

import (
	"os"

	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
)

var (
	GodogOptions = &godog.Options{
		Output: colors.Colored(os.Stdout),
		Format: "pretty",
	}
)
