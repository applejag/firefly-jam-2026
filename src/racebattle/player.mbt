///|
/// Rotation speed. Angle is the rotation speed per second
let rotation : Angle = Angle::from_deg(180.0 / 60.0)

///|
const MOVE_SPEED : Float = 0.5

///|
pub struct Player {
  // manually doing "struct of arrays" optimization,
  // mostly as an experiment to see how tedious it is,
  // but also because MoonBit doesn't allow "#valtype"s in other "#valtype"s
  peer : @firefly.Peer?
  mut pos : Vec2
  mut angle : Angle
}

///|
pub fn Player::new(peer : @firefly.Peer?) -> Player {
  Player::{ peer, pos: Vec2::default(), angle: Angle::default() }
}

///|
pub fn Player::update(self : Player) -> Unit {
  if self.peer is Some(peer) {
    let pad_maybe = @firefly.read_pad(peer~)
    if pad_maybe is Some(pad) {
      if pad.x != 0 && pad.y != 0 {
        let new_angle = @util.rotate_toward(self.angle, pad.azimuth(), rotation)
        self.angle = new_angle
      }
    }
  }
  let dir = Vec2::from_angle(self.angle).scale(MOVE_SPEED)
  self.pos = self.pos + dir.scale(MOVE_SPEED)
}

///|
pub fn Player::render(self : Player, world : World) -> Unit {
  let point = world.camera.world_vec2_to_camera_space(self.pos)
  let style = if Some(world.me) == self.peer {
    @firefly.Style::new_solid(DarkGreen)
  } else {
    @firefly.Style::new_solid(Red)
  }
  @firefly.draw_circle(point - Point::new(1, 1), 3, style)
  // draw arrow of movement direction
  let dir = Vec2::from_angle(self.angle)
  @firefly.draw_line(
    point + dir.scale(5).to_point(),
    point + dir.scale(8).to_point(),
    @firefly.LineStyle::new(Black),
  )
}
