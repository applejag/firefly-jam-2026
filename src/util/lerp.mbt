///|
pub fn lerp(from : Float, to : Float, weight : Float) -> Float {
  from + (to - from) * weight
}
