package camera

import (
	"fmt"
	"github.com/zgwit/iot-master/v4/lib"
	"github.com/zgwit/iot-master/v4/pkg/db"
)

var cameras lib.Map[Camera]

func Ensure(id string) (*Camera, error) {
	cam := cameras.Load(id)
	if cam == nil {
		err := Load(id)
		if err != nil {
			return nil, err
		}
		cam = cameras.Load(id)
	}
	return cam, nil
}

func Get(id string) *Camera {
	return cameras.Load(id)
}

func Set(id string, cam *Camera) {
	cameras.Store(id, cam)
}

func Load(id string) error {
	var cam Camera
	get, err := db.Engine.ID(id).Get(&cam)
	if err != nil {
		return err
	}
	if !get {
		return fmt.Errorf("camera %s not found", id)
	}

	return From(&cam)
}

func From(cam *Camera) error {
	cameras.Store(cam.Id, cam)
	return nil
}
