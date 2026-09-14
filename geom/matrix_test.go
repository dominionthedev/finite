package geom_test

import (
	"math"
	"testing"

	"github.com/dominionthedev/finite/geom"
)

func almostEqual(a, b float64) bool {
	return math.Abs(a-b) < 1e-9
}

func TestIdentity(t *testing.T) {
	m := geom.Identity()
	if !m.IsIdentity() {
		t.Fatal("Identity should report IsIdentity")
	}
	x, y := m.TransformPoint(3, 4)
	if !almostEqual(x, 3) || !almostEqual(y, 4) {
		t.Errorf("Identity transformed (3,4) to (%v,%v)", x, y)
	}
}

func TestTranslate(t *testing.T) {
	m := geom.Translate(10, -5)
	x, y := m.TransformPoint(0, 0)
	if !almostEqual(x, 10) || !almostEqual(y, -5) {
		t.Errorf("got (%v,%v)", x, y)
	}
	if m.String() != "translate(10,-5)" {
		t.Errorf("String = %q", m.String())
	}
}

func TestScale(t *testing.T) {
	m := geom.Scale(2, 3)
	x, y := m.TransformPoint(4, 5)
	if !almostEqual(x, 8) || !almostEqual(y, 15) {
		t.Errorf("got (%v,%v)", x, y)
	}
}

func TestScaleUniform(t *testing.T) {
	m := geom.ScaleUniform(2)
	x, y := m.TransformPoint(3, 4)
	if !almostEqual(x, 6) || !almostEqual(y, 8) {
		t.Errorf("got (%v,%v)", x, y)
	}
	if m.String() != "scale(2)" {
		t.Errorf("String = %q", m.String())
	}
}

func TestRotate90(t *testing.T) {
	m := geom.Rotate(90)
	x, y := m.TransformPoint(1, 0)
	if !almostEqual(x, 0) || !almostEqual(y, 1) {
		t.Errorf("rotate 90 of (1,0) = (%v,%v)", x, y)
	}
}

func TestRotateAround(t *testing.T) {
	// Rotate 180° around (10,10) should send (10,0) → (10,20)
	m := geom.RotateAround(180, 10, 10)
	x, y := m.TransformPoint(10, 0)
	if !almostEqual(x, 10) || !almostEqual(y, 20) {
		t.Errorf("got (%v,%v)", x, y)
	}
}

func TestSkewX(t *testing.T) {
	m := geom.SkewX(45) // tan(45)=1
	x, y := m.TransformPoint(0, 10)
	if !almostEqual(x, 10) || !almostEqual(y, 10) {
		t.Errorf("got (%v,%v)", x, y)
	}
}

func TestComposeOrder(t *testing.T) {
	// Compose applies matrices left-to-right via successive Mul.
	// With Mul defined as m.Mul(other) = apply other then m,
	// Compose(A, B) results in B applied first, then A when transforming points?
	// Verify and lock the observed (documented) behavior.
	m := geom.Compose(geom.Translate(10, 0), geom.ScaleUniform(2))
	x, y := m.TransformPoint(1, 1)
	// Observed: (12, 2) → scale first then translate.
	if !almostEqual(x, 12) || !almostEqual(y, 2) {
		t.Errorf("Compose(Translate, Scale) on (1,1): got (%v,%v), want (12,2)", x, y)
	}
}

func TestTransformBuilder(t *testing.T) {
	tr := geom.NewTransform().
		Translate(50, 50).
		Rotate(45).
		Translate(-50, -50)

	// Should be equivalent to RotateAround(45, 50, 50)
	expected := geom.RotateAround(45, 50, 50)
	m := tr.Matrix()

	// Compare by transforming a few points
	pts := [][2]float64{{0, 0}, {50, 50}, {100, 0}, {30, 70}}
	for _, p := range pts {
		x1, y1 := m.TransformPoint(p[0], p[1])
		x2, y2 := expected.TransformPoint(p[0], p[1])
		if !almostEqual(x1, x2) || !almostEqual(y1, y2) {
			t.Errorf("point (%v,%v): builder=(%v,%v) expected=(%v,%v)",
				p[0], p[1], x1, y1, x2, y2)
		}
	}
}

func TestMulAssociativity(t *testing.T) {
	a := geom.Translate(1, 2)
	b := geom.Scale(2, 3)
	c := geom.Rotate(30)

	left := a.Mul(b).Mul(c)
	right := a.Mul(b.Mul(c))

	pts := [][2]float64{{0, 0}, {1, 0}, {0, 1}, {5, 7}}
	for _, p := range pts {
		x1, y1 := left.TransformPoint(p[0], p[1])
		x2, y2 := right.TransformPoint(p[0], p[1])
		if !almostEqual(x1, x2) || !almostEqual(y1, y2) {
			t.Errorf("associativity failed at (%v,%v)", p[0], p[1])
		}
	}
}

func TestIsometric(t *testing.T) {
	m := geom.Isometric()
	if m.IsIdentity() {
		t.Fatal("Isometric should not be identity")
	}
	// Just ensure it produces a non-empty string and transforms something
	s := m.String()
	if s == "" {
		t.Error("Isometric String() should not be empty")
	}
	x, y := m.TransformPoint(10, 0)
	if almostEqual(x, 10) && almostEqual(y, 0) {
		t.Error("Isometric should change the point")
	}
}

func TestStringEmptyForIdentity(t *testing.T) {
	if geom.Identity().String() != "" {
		t.Error("Identity.String() should be empty")
	}
}
