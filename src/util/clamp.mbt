///|
fn[T : Compare] clamp(val : T, min : T, max : T) -> T {
  guard val > min else { min }
  guard val < max else { max }
  val
}
