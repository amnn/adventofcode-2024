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

func TestFlip(t *testing.T) {
	if DIR_U.Flip() != DIR_D {
		t.Errorf("expected %v, got %v", DIR_D, DIR_U.Flip())
	}

	if DIR_D.Flip() != DIR_U {
		t.Errorf("expected %v, got %v", DIR_U, DIR_D.Flip())
	}

	if DIR_L.Flip() != DIR_R {
		t.Errorf("expected %v, got %v", DIR_R, DIR_L.Flip())
	}

	if DIR_R.Flip() != DIR_L {
		t.Errorf("expected %v, got %v", DIR_L, DIR_R.Flip())
	}
}

func TestFlipMultiple(t *testing.T) {
	d := DIR_U | DIR_L
	if d.Flip() != DIR_D|DIR_R {
		t.Errorf("expected %v, got %v", DIR_D|DIR_R, d.Flip())
	}

	d = DIR_D | DIR_R
	if d.Flip() != DIR_U|DIR_L {
		t.Errorf("expected %v, got %v", DIR_U|DIR_L, d.Flip())
	}
}

func TestFlipSymmetric(t *testing.T) {
	d := DIR_U | DIR_D
	if d.Flip() != d {
		t.Errorf("expected %v, got %v", d, d.Flip())
	}

	d = DIR_L | DIR_R
	if d.Flip() != d {
		t.Errorf("expected %v, got %v", d, d.Flip())
	}
}
