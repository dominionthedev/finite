package draw

import (
	"fmt"
	"strings"
)

// Point is a 2D coordinate in user space.
type Point struct {
	X, Y float64
}

// Pt is a convenience constructor for Point.
func Pt(x, y float64) Point { return Point{X: x, Y: y} }

// Polygon is a closed shape defined by a list of points (<polygon>).
type Polygon struct {
	points []Point
	style  shapeStyle
}

// NewPolygon creates a polygon from the given vertices.
func NewPolygon(points []Point, opts ...ShapeOption) *Polygon {
	p := &Polygon{points: append([]Point(nil), points...)}
	p.style.opacity = 1.0
	for _, opt := range opts {
		opt(&p.style)
	}
	return p
}

// Render generates the SVG for the polygon.
func (p *Polygon) Render() (string, error) {
	attrs := fmt.Sprintf(` points="%s"%s`, formatPoints(p.points), p.style.attrString())
	el := elementMarkup("polygon", attrs, p.style.title, p.style.desc)
	return wrap(p.style.buildDefs(), el), nil
}

// Polyline is an open path of connected segments (<polyline>).
type Polyline struct {
	points []Point
	style  shapeStyle
}

// NewPolyline creates a polyline from the given vertices.
// Default fill is none so only the stroke is visible unless Fill is set.
func NewPolyline(points []Point, opts ...ShapeOption) *Polyline {
	p := &Polyline{points: append([]Point(nil), points...)}
	p.style.opacity = 1.0
	p.style.fill = "none"
	for _, opt := range opts {
		opt(&p.style)
	}
	return p
}

// Render generates the SVG for the polyline.
func (p *Polyline) Render() (string, error) {
	attrs := fmt.Sprintf(` points="%s"%s`, formatPoints(p.points), p.style.attrString())
	el := elementMarkup("polyline", attrs, p.style.title, p.style.desc)
	return wrap(p.style.buildDefs(), el), nil
}

func formatPoints(pts []Point) string {
	var b strings.Builder
	for i, pt := range pts {
		if i > 0 {
			b.WriteByte(' ')
		}
		b.WriteString(fmt.Sprintf("%.4f,%.4f", pt.X, pt.Y))
	}
	return b.String()
}
