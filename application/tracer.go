package application

import "fmt"

const (
	ScopeName = "github.com/yoshino-s/go-framework/application"
)

func (a *EmptyApplication) spanName(stage ApplicationStage) string {
	switch stage {
	case StageBeforeSetup:
		return fmt.Sprintf("%s.BeforeSetup", a.Name)
	case StageSetup:
		return fmt.Sprintf("%s.Setup", a.Name)
	case StageAfterSetup:
		return fmt.Sprintf("%s.AfterSetup", a.Name)
	case StageRun:
		return fmt.Sprintf("%s.Run", a.Name)
	case StageClose:
		return fmt.Sprintf("%s.Close", a.Name)
	}
	return "unknown"
}
