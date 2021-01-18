package env

import "fmt"

type Environment string

const (
	EnvironmentProd         Environment = "prod"
	EnvironmentStage        Environment = "stage"
	EnvironmentLocal        Environment = "local"
	EnvironmentLocalNoEmail Environment = "local-no-emails"
)

func (e Environment) Str() string {
	return string(e)
}

func mustEnvironmentFromString(s string) Environment {
	switch s {
	case "prod":
		return EnvironmentProd
	case "stage":
		return EnvironmentStage
	case "local":
		return EnvironmentLocal
	case "local-no-email":
		return EnvironmentLocalNoEmail
	default:
		panic(fmt.Sprintf("Unrecognized environment: %s", s))
	}
}

func GetAbsoluteURLForEnvironment(path string) string {
	environment := mustEnvironmentFromString(GetEnvironmentVariableOrDefault("ENV", EnvironmentProd.Str()))
	switch environment {
	case EnvironmentProd:
		return fmt.Sprintf("https://www.babblegraph.com/%s", path)
	case EnvironmentStage:
		panic("unimplemented")
	case EnvironmentLocal,
		EnvironmentLocalNoEmail:
		return fmt.Sprintf("http://localhost:8080/%s", path)
	default:
		panic(fmt.Sprintf("Unrecognized environment: %s", environment.Str()))
	}

}
