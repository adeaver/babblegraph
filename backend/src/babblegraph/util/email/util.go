package email

import (
	"babblegraph/util/env"
	"babblegraph/util/ptr"
	"fmt"
	"os"
)

func getPathForTemplateFile(filename string) (*string, error) {
	templatePath := env.GetEnvironmentVariableOrDefault("TEMPLATES_PATH", "/util/email/templates/")
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	return ptr.String(fmt.Sprintf("%s%s%s", cwd, templatePath, filename)), nil
}
