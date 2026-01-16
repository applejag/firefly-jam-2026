///|
using @firefly {type Point}

///|
#valtype
pub struct Vec2 {
  x : Float
  y : Float
} derive(Show, Eq, Compare, Hash, Default)

///|
#inline
pub fn Vec2::new(x : Float, y : Float) -> Vec2 {
  Vec2::{ x, y }
}

///|
#inline
pub fn Vec2::from_point(point : Point) -> Vec2 {
  Vec2::{ x: Float::from_int(point.x), y: Float::from_int(point.y) }
}

///|
#inline
pub fn Vec2::to_point(self : Vec2) -> Point {
  Point::new(self.x.to_int(), self.y.to_int())
}

///|
#inline
pub fn Vec2::from_angle(angle : Angle) -> Vec2 {
  Vec2::{
    x: @micromath.cos(angle.to_rad()),
    y: -@micromath.sin(angle.to_rad()),
  }
}

///|
#inline
pub fn Vec2::round(self : Vec2) -> Vec2 {
  Vec2::{ x: self.x.round(), y: self.y.round() }
}

///|
/// Returns a `Vec2` with absolute (non-negative) `x` and `y` values.
pub fn Vec2::abs(self : Vec2) -> Vec2 {
  Vec2::{ x: self.x.abs(), y: self.y.abs() }
}

///|
/// Add together two Vec2s.
pub impl Add for Vec2 with add(self : Vec2, other : Vec2) -> Vec2 {
  { x: self.x + other.x, y: self.y + other.y }
}

///|
/// Subtract one Vec2 from another.
pub impl Sub for Vec2 with sub(self : Vec2, other : Vec2) -> Vec2 {
  { x: self.x - other.x, y: self.y - other.y }
}

///|
/// Negate a Vec2.
pub impl Neg for Vec2 with neg(self : Vec2) -> Vec2 {
  { x: -self.x, y: -self.y }
}

///|
/// Scale the `x` and `y` of a Vec2 by a scalar.
pub fn Vec2::scale(self : Vec2, factor : Float) -> Vec2 {
  { x: self.x * factor, y: self.y * factor }
}
