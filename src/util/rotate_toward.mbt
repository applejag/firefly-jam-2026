///|
using @firefly {type Angle}

///|
pub fn rotate_toward(from : Angle, to : Angle, delta : Angle) -> Angle {
  let difference = from.difference(to).to_rad()
  let abs_difference = difference.abs()
  Angle::from_rad(
    from.to_rad() +
    clamp(delta.to_rad(), abs_difference - @micromath.PI, abs_difference) *
    @micromath.signum(difference),
  )
}
