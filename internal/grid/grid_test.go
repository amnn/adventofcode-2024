package grid

import (
	"testing"
)

func TestRotateClockwise(t *testing.T) {
	d := DIR_U

	for _, dir := range []Dir{DIR_R, DIR_D, DIR_L, DIR_U} {
		if d = d.RotateClockwise(); d != dir {
			t.Errorf("expected %v, got %v", dir, d)
		}
	}
}

func TestRotateCounterClockwise(t *testing.T) {
	d := DIR_U
	for _, dir := range []Dir{DIR_L, DIR_D, DIR_R, DIR_U} {
		if d = d.RotateCounterClockwise(); d != dir {
			t.Errorf("expected %v, got %v", dir, d)
		}
	}
}

func TestRotateClockwiseMultiple(t *testing.T) {
	d := DIR_U | DIR_L
	for _, dir := range []Dir{DIR_U | DIR_R, DIR_R | DIR_D, DIR_D | DIR_L, DIR_L | DIR_U} {
		if d = d.RotateClockwise(); d != dir {
			t.Errorf("expected %v, got %v", dir, d)
		}
	}
}

func TestRotateCounterClockwiseMultiple(t *testing.T) {
	d := DIR_U | DIR_L
	for _, dir := range []Dir{DIR_L | DIR_D, DIR_D | DIR_R, DIR_R | DIR_U, DIR_U | DIR_L} {
		if d = d.RotateCounterClockwise(); d != dir {
			t.Errorf("expected %v, got %v", dir, d)
		}
	}
}
