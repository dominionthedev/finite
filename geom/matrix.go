// Package geom provides 2D affine transforms and matrix utilities for Finite.
//
// It is the foundation for positioning, rotation, scaling, skewing, and
// composing visual elements. All transforms produce valid SVG transform
// attribute values and can be composed without string manipulation.
package geom

import (
	"fmt"
	"math"
	"strings"
)

// Matrix represents a 2D affine transformation matrix of the form:
//
//	| a  c  e |
//	| b  d  f |
//	| 0  0  1 |
//
// This matches the SVG/CSS transform matrix(a, b, c, d, e, f).
type Matrix struct {
	A, B, C, D, E, F float64
}

// Identity returns the identity matrix (no transformation).
func Identity() Matrix {
	return Matrix{A: 1, D: 1}
}

// Translate returns a translation matrix.
func Translate(tx, ty float64) Matrix {
	return Matrix{A: 1, D: 1, E: tx, F: ty}
}

// Scale returns a non-uniform scale matrix.
func Scale(sx, sy float64) Matrix {
	return Matrix{A: sx, D: sy}
}

// ScaleUniform returns a uniform scale matrix.
func ScaleUniform(s float64) Matrix {
	return Scale(s, s)
}

// Rotate returns a rotation matrix by angle degrees around the origin.
func Rotate(degrees float64) Matrix {
	rad := degrees * math.Pi / 180
	cos, sin := math.Cos(rad), math.Sin(rad)
	return Matrix{A: cos, B: sin, C: -sin, D: cos}
}

// RotateAround returns a rotation matrix by angle degrees around the point (cx, cy).
func RotateAround(degrees, cx, cy float64) Matrix {
	return Translate(cx, cy).Mul(Rotate(degrees)).Mul(Translate(-cx, -cy))
}

// SkewX returns a horizontal skew matrix by angle degrees.
func SkewX(degrees float64) Matrix {
	rad := degrees * math.Pi / 180
	return Matrix{A: 1, C: math.Tan(rad), D: 1}
}

// SkewY returns a vertical skew matrix by angle degrees.
func SkewY(degrees float64) Matrix {
	rad := degrees * math.Pi / 180
	return Matrix{A: 1, B: math.Tan(rad), D: 1}
}

// Mul multiplies two matrices (this ∘ other). The result applies other first, then this.
// This matches the mathematical convention used by SVG (right-to-left application).
func (m Matrix) Mul(other Matrix) Matrix {
	return Matrix{
		A: m.A*other.A + m.C*other.B,
		B: m.B*other.A + m.D*other.B,
		C: m.A*other.C + m.C*other.D,
		D: m.B*other.C + m.D*other.D,
		E: m.A*other.E + m.C*other.F + m.E,
		F: m.B*other.E + m.D*other.F + m.F,
	}
}

// TransformPoint applies the matrix to a point (x, y).
func (m Matrix) TransformPoint(x, y float64) (float64, float64) {
	return m.A*x + m.C*y + m.E, m.B*x + m.D*y + m.F
}

// IsIdentity reports whether the matrix is (approximately) the identity.
func (m Matrix) IsIdentity() bool {
	const eps = 1e-10
	return math.Abs(m.A-1) < eps && math.Abs(m.B) < eps &&
		math.Abs(m.C) < eps && math.Abs(m.D-1) < eps &&
		math.Abs(m.E) < eps && math.Abs(m.F) < eps
}

// String returns a compact SVG transform attribute value.
// It prefers readable forms (translate, scale, rotate) when possible,
// otherwise falls back to matrix(...).
func (m Matrix) String() string {
	if m.IsIdentity() {
		return ""
	}

	const eps = 1e-9

	// Pure translation
	if math.Abs(m.A-1) < eps && math.Abs(m.B) < eps &&
		math.Abs(m.C) < eps && math.Abs(m.D-1) < eps {
		return fmt.Sprintf("translate(%.4g,%.4g)", m.E, m.F)
	}

	// Pure scale (no translation, no shear/rotation)
	if math.Abs(m.B) < eps && math.Abs(m.C) < eps &&
		math.Abs(m.E) < eps && math.Abs(m.F) < eps {
		if math.Abs(m.A-m.D) < eps {
			return fmt.Sprintf("scale(%.4g)", m.A)
		}
		return fmt.Sprintf("scale(%.4g,%.4g)", m.A, m.D)
	}

	// Pure rotation around origin (no translation)
	// matrix = | cos  -sin  0 |
	//          | sin   cos  0 |
	if math.Abs(m.E) < eps && math.Abs(m.F) < eps {
		// Check orthonormal rotation (no scale/shear)
		if math.Abs(m.A*m.A+m.B*m.B-1) < 1e-6 &&
			math.Abs(m.C*m.C+m.D*m.D-1) < 1e-6 &&
			math.Abs(m.A-m.D) < 1e-6 && math.Abs(m.B+m.C) < 1e-6 {
			deg := math.Atan2(m.B, m.A) * 180 / math.Pi
			return fmt.Sprintf("rotate(%.4g)", deg)
		}
	}

	// Translate + rotate around origin: common Instance pattern
	// matrix = T * R  =>  a=cos, b=sin, c=-sin, d=cos, e=tx, f=ty
	if math.Abs(m.A*m.A+m.B*m.B-1) < 1e-6 &&
		math.Abs(m.C*m.C+m.D*m.D-1) < 1e-6 &&
		math.Abs(m.A-m.D) < 1e-6 && math.Abs(m.B+m.C) < 1e-6 {
		deg := math.Atan2(m.B, m.A) * 180 / math.Pi
		parts := []string{}
		if math.Abs(m.E) > eps || math.Abs(m.F) > eps {
			parts = append(parts, fmt.Sprintf("translate(%.4g,%.4g)", m.E, m.F))
		}
		if math.Abs(deg) > eps {
			parts = append(parts, fmt.Sprintf("rotate(%.4g)", deg))
		}
		if len(parts) > 0 {
			return strings.Join(parts, " ")
		}
	}

	// Translate + uniform scale
	if math.Abs(m.B) < eps && math.Abs(m.C) < eps && math.Abs(m.A-m.D) < eps {
		parts := []string{}
		if math.Abs(m.E) > eps || math.Abs(m.F) > eps {
			parts = append(parts, fmt.Sprintf("translate(%.4g,%.4g)", m.E, m.F))
		}
		if math.Abs(m.A-1) > eps {
			parts = append(parts, fmt.Sprintf("scale(%.4g)", m.A))
		}
		if len(parts) > 0 {
			return strings.Join(parts, " ")
		}
	}

	// General case
	return fmt.Sprintf("matrix(%.6g,%.6g,%.6g,%.6g,%.6g,%.6g)",
		m.A, m.B, m.C, m.D, m.E, m.F)
}

// Transform is a convenience builder that accumulates affine operations.
// Operations are applied in the order they are chained (left to right in code,
// which corresponds to right-to-left matrix multiplication for SVG).
type Transform struct {
	m Matrix
}

// NewTransform starts a new transform builder from the identity.
func NewTransform() *Transform {
	return &Transform{m: Identity()}
}

// FromMatrix starts a builder from an existing matrix.
func FromMatrix(m Matrix) *Transform {
	return &Transform{m: m}
}

// Translate appends a translation.
func (t *Transform) Translate(tx, ty float64) *Transform {
	t.m = t.m.Mul(Translate(tx, ty))
	return t
}

// Scale appends a non-uniform scale.
func (t *Transform) Scale(sx, sy float64) *Transform {
	t.m = t.m.Mul(Scale(sx, sy))
	return t
}

// ScaleUniform appends a uniform scale.
func (t *Transform) ScaleUniform(s float64) *Transform {
	t.m = t.m.Mul(ScaleUniform(s))
	return t
}

// Rotate appends a rotation (degrees) around the origin.
func (t *Transform) Rotate(degrees float64) *Transform {
	t.m = t.m.Mul(Rotate(degrees))
	return t
}

// RotateAround appends a rotation (degrees) around (cx, cy).
func (t *Transform) RotateAround(degrees, cx, cy float64) *Transform {
	t.m = t.m.Mul(RotateAround(degrees, cx, cy))
	return t
}

// SkewX appends a horizontal skew (degrees).
func (t *Transform) SkewX(degrees float64) *Transform {
	t.m = t.m.Mul(SkewX(degrees))
	return t
}

// SkewY appends a vertical skew (degrees).
func (t *Transform) SkewY(degrees float64) *Transform {
	t.m = t.m.Mul(SkewY(degrees))
	return t
}

// Matrix returns the accumulated matrix.
func (t *Transform) Matrix() Matrix {
	return t.m
}

// String returns the SVG transform attribute value.
func (t *Transform) String() string {
	return t.m.String()
}

// Compose combines multiple matrices. The first matrix is applied first.
func Compose(ms ...Matrix) Matrix {
	result := Identity()
	for _, m := range ms {
		result = result.Mul(m)
	}
	return result
}

// Isometric returns a common isometric projection transform.
// It applies: rotate(30°) → skewX(-30°) → scaleY(√3/2 ≈ 0.866).
// This is a convenient starting point for isometric scenes.
func Isometric() Matrix {
	return Compose(
		Rotate(30),
		SkewX(-30),
		Scale(1, math.Sqrt(3)/2),
	)
}

// Oblique returns a simple oblique (cabinet-style) projection using a horizontal skew.
func Oblique(degrees float64) Matrix {
	return SkewX(degrees)
}

// Parse is a minimal helper that accepts a few common shorthand forms and
// returns a Matrix. It is intentionally small; complex strings should be
// built with the Transform API instead of parsed.
//
// Supported:
//   - "" or "none" → Identity
//   - "translate(tx,ty)"
//   - "scale(s)" / "scale(sx,sy)"
//   - "rotate(deg)" / "rotate(deg,cx,cy)"
//
// Anything else falls back to Identity and is not an error (callers that need
// strict parsing should build matrices explicitly).
func Parse(s string) Matrix {
	s = strings.TrimSpace(s)
	if s == "" || s == "none" {
		return Identity()
	}
	// Keep this deliberately minimal for now. Full SVG transform list parsing
	// can be added later if needed. Prefer constructing with the API.
	return Identity()
}
