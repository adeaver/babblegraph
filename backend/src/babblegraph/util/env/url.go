package env

import "fmt"

func GetAbsoluteURLForEnvironment(path string) string {
	environment := mustEnvironmentFromString(GetEnvironmentVariableOrDefault("ENV", EnvironmentProd.Str()))
	switch environment {
	case EnvironmentProd:
		return fmt.Sprintf("https://www.babblegraph.com/%s", path)
	case EnvironmentStage:
		panic("unimplemented")
	case EnvironmentLocal,
		EnvironmentLocalTestEmail,
		EnvironmentLocalNoEmail:
		return fmt.Sprintf("http://localhost:8080/%s", path)
	default:
		panic(fmt.Sprintf("Unrecognized environment: %s", environment.Str()))
	}

}
