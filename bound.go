package orb

import (
	"math"
)

// A Bound represents an enclosed "box" on the sphere.
// It does not know anything about the anti-meridian (TODO).
type Bound struct {
	Min, Max Point
}

// NewBoundFromPoints creates a new bound given two opposite corners.
// These corners can be either sw/ne or se/nw.
func NewBoundFromPoints(corner, oppositeCorner Point) Bound {
	return Bound{corner, corner}.Extend(oppositeCorner)
}

// GeoJSONType returns the GeoJSON type for the object.
func (b Bound) GeoJSONType() string {
	return "Polygon"
}

// Dimensions returns 2 because a Bound is a 2d object.
func (b Bound) Dimensions() int {
	return 2
}

// ToPolygon converts the bound into a Polygon object.
func (b Bound) ToPolygon() Polygon {
	return Polygon{b.ToRing()}
}

// ToRing converts the bound into a loop defined
// by the boundary of the box.
func (b Bound) ToRing() Ring {
	return Ring{
		b.Min,
		Point{b.Max[0], b.Min[1]},
		b.Max,
		Point{b.Min[0], b.Max[1]},
		b.Min,
	}
}

// Extend grows the bound to include the new point.
func (b Bound) Extend(point Point) Bound {
	// already included, no big deal
	if b.Contains(point) {
		return b
	}

	return Bound{
		Min: Point{
			math.Min(b.Min[0], point[0]),
			math.Min(b.Min[1], point[1]),
		},
		Max: Point{
			math.Max(b.Max[0], point[0]),
			math.Max(b.Max[1], point[1]),
		},
	}
}

// Union extends this bound to contain the union of this and the given bound.
func (b Bound) Union(other Bound) Bound {
	b = b.Extend(other.Min)
	b = b.Extend(other.Max)
	b = b.Extend(other.LeftTop())
	b = b.Extend(other.RightBottom())

	return b
}

// Contains determines if the point is within the bound.
// Points on the boundary are considered within.
func (b Bound) Contains(point Point) bool {
	if point[1] < b.Min[1] || b.Max[1] < point[1] {
		return false
	}

	if point[0] < b.Min[0] || b.Max[0] < point[0] {
		return false
	}

	return true
}

// Intersects determines if two bounds intersect.
// Returns true if they are touching.
func (b Bound) Intersects(bound Bound) bool {
	if (b.Max[0] < bound.Min[0]) ||
		(b.Min[0] > bound.Max[0]) ||
		(b.Max[1] < bound.Min[1]) ||
		(b.Min[1] > bound.Max[1]) {
		return false
	}

	return true
}

// Center returns the center of the bounds by "averaging" the x and y coords.
func (b Bound) Center() Point {
	return Point{
		(b.Min[0] + b.Max[0]) / 2.0,
		(b.Min[1] + b.Max[1]) / 2.0,
	}
}

// Top returns the top of the bound.
func (b Bound) Top() float64 {
	return b.Max[1]
}

// Bottom returns the bottom of the bound.
func (b Bound) Bottom() float64 {
	return b.Min[1]
}

// Right returns the right of the bound.
func (b Bound) Right() float64 {
	return b.Max[0]
}

// Left returns the left of the bound.
func (b Bound) Left() float64 {
	return b.Min[0]
}

// LeftTop returns the upper left point of the bound.
func (b Bound) LeftTop() Point {
	return Point{b.Left(), b.Top()}
}

// RightBottom return the lower right point of the bound.
func (b Bound) RightBottom() Point {
	return Point{b.Right(), b.Bottom()}
}

// IsEmpty returns true if it contains zero area or if
// it's in some malformed negative state where the left point is larger than the right.
// This can be caused by padding too much negative.
func (b Bound) IsEmpty() bool {
	return b.Min[0] > b.Max[0] || b.Min[1] > b.Max[1]
}

// IsZero return true if the bound just includes just null island.
func (b Bound) IsZero() bool {
	return b == Bound{}
}

// Bound returns the the same bound.
func (b Bound) Bound() Bound {
	return b
}

// Equal returns if two bounds are equal.
func (b Bound) Equal(c Bound) bool {
	return b.Min == c.Min && b.Max == c.Max
}