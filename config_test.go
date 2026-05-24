package qrgode

import "testing"

func TestNewSolidColor(t *testing.T) {
	c := NewSolidColor("#abcdef")
	if c.Type() != "solid" {
		t.Errorf("Type = %s", c.Type())
	}
}

func TestNewLinearGradientColor(t *testing.T) {
	c := NewLinearGradientColor(45, []string{"#000", "#fff"})
	if c.Type() != "linear-gradient" {
		t.Errorf("Type = %s", c.Type())
	}
}

func TestNewRadialGradientColor(t *testing.T) {
	c := NewRadialGradientColor(0.5, 0.5, []string{"#000", "#fff"})
	if c.Type() != "radial-gradient" {
		t.Errorf("Type = %s", c.Type())
	}
}
