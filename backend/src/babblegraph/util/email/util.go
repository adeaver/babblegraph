package email

import (
	"babblegraph/util/ptr"
	"fmt"
	"os"
)

const templatePath = "/util/email/templates/"

func getPathForTemplateFile(filename string) (*string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	return ptr.String(fmt.Sprintf("%s%s%s", cwd, templatePath, filename)), nil
}
