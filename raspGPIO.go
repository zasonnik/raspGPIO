package raspGPIO

import (
	"fmt"
	//    "raspGPIO"
	"os"
	"time"
)

const (
	Input = iota
	Output
)

const (
	Low = iota
	Hight
)

type Pin struct {
	Pin_id uint8
}

func (p *Pin) Unexport() {
	file, _ := os.OpenFile("/sys/class/gpio/unexport", os.O_WRONLY|os.O_SYNC, os.ModeDevice|os.ModeCharDevice)
	defer file.Close()
	_, _ = fmt.Fprintf(file, "%d\n", p.Pin_id)
}

func (p *Pin) Export() {
	file, err := os.OpenFile("/sys/class/gpio/export", os.O_WRONLY|os.O_SYNC, os.ModeDevice|os.ModeCharDevice)
	for err != nil {
		panic("NO GPIO")
	}
	defer file.Close()
	_, _ = fmt.Fprintf(file, "%d\n", p.Pin_id)
	//	time.Sleep(50 * time.Millisecond)
}

func (p *Pin) SetDirection(direction int) error {
	file, err := os.OpenFile(fmt.Sprintf("/sys/class/gpio/gpio%d/direction", p.Pin_id), os.O_WRONLY|os.O_SYNC, os.ModeDevice|os.ModeCharDevice)
	for err != nil {
		file, err = os.OpenFile(fmt.Sprintf("/sys/class/gpio/gpio%d/direction", p.Pin_id), os.O_WRONLY|os.O_SYNC, os.ModeDevice|os.ModeCharDevice)
	}
	defer file.Close()
	if direction != Input {
		fmt.Fprintln(file, "out")
	} else {
		fmt.Fprintln(file, "in")
	}
	return nil
}

func (p *Pin) ReadValue() (int, error) {
	file, err := os.OpenFile(fmt.Sprintf("/sys/class/gpio/gpio%d/value", p.Pin_id), os.O_RDONLY|os.O_SYNC, os.ModeDevice|os.ModeCharDevice)
	if err != nil {
		return Low, err
	}
	defer file.Close()
	data := make([]byte, 1)
	_, _ = file.Read(data)
	if string(data[:]) == "1" {
		return Hight, nil
	} else {
		return Low, nil
	}
}

func (p *Pin) WriteValue(value int) error {
	file, err := os.OpenFile(fmt.Sprintf("/sys/class/gpio/gpio%d/value", p.Pin_id), os.O_WRONLY|os.O_SYNC, os.ModeDevice|os.ModeCharDevice)
	if err != nil {
		return err
	}
	defer file.Close()
	if value == Low {
		fmt.Fprint(file, "1")
	} else {
		fmt.Fprint(file, "0")
	}
	return nil
}

func (p *Pin) ReadValueToChan(c chan int) {
	file, _ := os.OpenFile(fmt.Sprintf("/sys/class/gpio/gpio%d/value", p.Pin_id), os.O_RDONLY|os.O_SYNC, os.ModeDevice|os.ModeCharDevice)
	data := make([]byte, 1)
	lastValue := Low
	nowValue := Low
	for {
		_, _ = file.ReadAt(data, 0)
		if string(data[:]) == "1" {
			nowValue = Hight
		} else {
			nowValue = Low
		}
		if nowValue != lastValue {
			c <- nowValue
			time.Sleep(100 * time.Millisecond)
		}
		lastValue = nowValue
	}
}