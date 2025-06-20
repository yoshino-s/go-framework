package application

import "fmt"

const (
	ScopeName = "github.com/yoshino-s/go-framework/application"
)

func (a *EmptyApplication) spanName(stage ApplicationStage) string {
	switch stage {
	case StageInitialize:
		return fmt.Sprintf("%s.Initialize", a.Name)
	case StageSetup:
		return fmt.Sprintf("%s.Setup", a.Name)
	case StageRun:
		return fmt.Sprintf("%s.Run", a.Name)
	case StageClose:
		return fmt.Sprintf("%s.Close", a.Name)
	}
	return "unknown"
}
