package blog

import (
	"babblegraph/util/env"
	"fmt"
)

func MakeImageDirectory(blogID ID) string {
	envName := env.MustEnvironmentName()
	return fmt.Sprintf("blog-content/%s/images/%s", envName.Str(), blogID)
}
