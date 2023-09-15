package testutils

import (
	"os"

	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
)

func GodogOptions(featureFilePaths ...string) (options *godog.Options) {
	options = &godog.Options{
		Output: colors.Colored(os.Stdout),
		Format: "pretty",
	}

	options.Paths = featureFilePaths

	return
}
