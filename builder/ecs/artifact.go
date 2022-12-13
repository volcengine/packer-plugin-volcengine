package ecs

import (
	"fmt"

	"github.com/volcengine/volcengine-go-sdk/service/ecs"
	"github.com/volcengine/volcengine-go-sdk/volcengine"
)

type Artifact struct {
	VolcengineImageId string
	BuilderIdValue    string
	Client            *VolcengineClientWrapper
}

func (a *Artifact) BuilderId() string {
	return a.BuilderIdValue
}

func (a *Artifact) Files() []string {
	return nil
}

func (a *Artifact) Id() string {
	return a.VolcengineImageId
}

func (a *Artifact) String() string {
	return fmt.Sprintf("Volcengine images were created:%s", a.VolcengineImageId)
}

func (a *Artifact) State(name string) interface{} {
	return nil
}

func (a *Artifact) Destroy() error {
	req := ecs.DeleteImagesInput{
		ImageIds: volcengine.StringSlice([]string{a.VolcengineImageId}),
	}
	_, err := a.Client.EcsClient.DeleteImages(&req)
	if err != nil {
		return err
	}
	return nil
}
