///|
pub fn move_towards(start : Float, end : Float, step : Float) -> Float {
  if (end - start).abs() <= step {
    end
  } else {
    start + @micromath.signum(end - start) * step
  }
}

///|
test "move_towards no change" {
  assert_eq(move_towards(5, 5, 1), 5)
  assert_eq(move_towards(5, 10, 0), 5)
}

///|
test "move_towards positive step" {
  assert_eq(move_towards(5, 10, 1), 6)
  assert_eq(move_towards(10, 5, 1), 9)
}

///|
test "move_towards negative step" {
  assert_eq(move_towards(5, 10, -1), 4)
  assert_eq(move_towards(10, 5, -1), 11)
}

///|
test "move_towards overstep" {
  assert_eq(move_towards(5, 10, 10), 10)
  assert_eq(move_towards(10, 5, 10), 5)
}
