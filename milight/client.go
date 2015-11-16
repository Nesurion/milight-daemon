package milight

import (
	"fmt"
	"runtime"
	"strings"
	"time"

	"github.com/lucasb-eyer/go-colorful"
	"github.com/nesurion/go-limitless"
)

type MilightClient struct {
	Controller *limitless.LimitlessController
}

func NewClient(conf Config) (*MilightClient, error) {
	c, err := limitless.NewLimitlessController(conf.Bridge)
	if err != nil {
		return nil, err
	}
	m := MilightClient{
		Controller: c,
	}
	setGroups(m.Controller)
	return &m, nil
}

func (m *MilightClient) On(id int) error {
	if id == -1 {
		err := m.Controller.AllOn()
		if err != nil {
			err = fmt.Errorf("failed to send %s", funcName())
		}
		return err
	}

	err := m.Controller.Groups[id].On()
	if err != nil {
		err := fmt.Errorf("failed to send %s", funcName())
		return err
	}
	return nil
}

func (m *MilightClient) Off(id int) error {
	if id == -1 {
		err := m.Controller.AllOff()
		if err != nil {
			err = fmt.Errorf("failed to send %s", funcName())
		}
		return err
	}

	err := m.Controller.Groups[id].Off()
	if err != nil {
		err = fmt.Errorf("failed to send %s", funcName())
	}
	return err
}

func (m *MilightClient) Rgb(id int, color colorful.Color) error {
	if id == -1 {
		for _, g := range m.Controller.Groups {
			err := g.SendColor(color)
			if err != nil {
				err := fmt.Errorf("failed to send %s", funcName())
				return err
			}
			wait()
		}
		return nil
	}

	err := m.Controller.Groups[id].SendColor(color)
	if err != nil {
		err := fmt.Errorf("failed to send %s", funcName())
		return err
	}
	return nil
}

func (m *MilightClient) Brightness(id int, brightness uint8) error {
	if id == -1 {
		for _, g := range m.Controller.Groups {
			err := g.SetBri(brightness)
			if err != nil {
				err := fmt.Errorf("failed to send %s", funcName())
				return err
			}
			wait()
		}
		return nil
	}

	err := m.Controller.Groups[id].SetBri(brightness)
	if err != nil {
		err := fmt.Errorf("failed to send %s", funcName())
		return err
	}
	return nil
}

func (m *MilightClient) Color(id int, color uint8) error {
	if id == -1 {
		for _, g := range m.Controller.Groups {
			err := g.SetHue(color)
			if err != nil {
				err := fmt.Errorf("failed to send %s", funcName())
				return err
			}
			wait()
		}
		return nil
	}

	err := m.Controller.Groups[id].SetHue(color)
	if err != nil {
		err := fmt.Errorf("failed to send %s", funcName())
		return err
	}
	return nil
}

func (m *MilightClient) White(id int) error {
	if id == -1 {
		for _, g := range m.Controller.Groups {
			err := g.White()
			if err != nil {
				err := fmt.Errorf("failed to send %s", funcName())
				return err
			}
			wait()
		}
		return nil
	}

	err := m.Controller.Groups[id].White()
	if err != nil {
		err := fmt.Errorf("failed to send %s", funcName())
		return err
	}
	return nil
}

func (m *MilightClient) Night(id int) error {
	if id == -1 {
		for _, g := range m.Controller.Groups {
			err := g.Night()
			if err != nil {
				err := fmt.Errorf("failed to send %s", funcName())
				return err
			}
			wait()
		}
		return nil
	}

	err := m.Controller.Groups[id].Night()
	if err != nil {
		err := fmt.Errorf("failed to send %s", funcName())
		return err
	}
	return nil
}

func (m *MilightClient) Disco(id int, speed string) error {
	switch speed {
	case "up":
		if id == -1 {
			for _, g := range m.Controller.Groups {
				err := g.DiscoFaster()
				if err != nil {
					err := fmt.Errorf("failed to send %s", funcName())
					return err
				}
				wait()
			}
			return nil
		}

		err := m.Controller.Groups[id].DiscoFaster()
		if err != nil {
			err := fmt.Errorf("failed to send %s", funcName())
			return err
		}
		return nil
	case "down":
		if id == -1 {
			for _, g := range m.Controller.Groups {
				err := g.DiscoSlower()
				if err != nil {
					err := fmt.Errorf("failed to send %s", funcName())
					return err
				}
				wait()
			}
			return nil
		}

		err := m.Controller.Groups[id].DiscoSlower()
		if err != nil {
			err := fmt.Errorf("failed to send %s", funcName())
			return err
		}
		return nil
	default:
		if id == -1 {
			for _, g := range m.Controller.Groups {
				err := g.Disco()
				if err != nil {
					err := fmt.Errorf("failed to send %s", funcName())
					return err
				}
				wait()
			}
			return nil
		}

		err := m.Controller.Groups[id].Disco()
		if err != nil {
			err := fmt.Errorf("failed to send %s", funcName())
			return err
		}
		return nil
	}
}

func setGroups(c *limitless.LimitlessController) {
	g := make([]limitless.LimitlessGroup, 4, 4)
	for i := 0; i < 4; i++ {
		g[i] = limitless.LimitlessGroup{
			Id:         i + 1,
			Controller: c,
		}
	}
	c.Groups = g
}

func funcName() string {
	pc, _, _, _ := runtime.Caller(1)
	t := strings.Split(runtime.FuncForPC(pc).Name(), ".")
	return strings.ToLower(t[len(t)-1])
}

func wait() {
	time.Sleep(200 * time.Millisecond)
}
