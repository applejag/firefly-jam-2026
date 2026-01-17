///|
const FPS : Float = 60

///|
/// Rotation speed at top speed. Angle is the rotation speed per second
let rotation_speed_when_moving : Angle = Angle::from_deg(110.0 / FPS)

///|
/// Rotation speed when standing completely still. Angle is the rotation speed per second
let rotation_speed_when_still : Angle = Angle::from_deg(540.0 / FPS)

///|
/// Maximum movement speed, in pixels/frame
const MOVE_MAX_SPEED : Float = 1.2

///|
/// Acceleration from 0 to max speed, in %/frame
const MOVE_ACCELERATION : Float = 0.00370370370 // 0%-100% in 1/270th, i.e in 4.5sec

///|
/// Deacceleration (break force) to go from max speed to 0, in %/frame
const MOVE_DEACCELERATION : Float = 0.016666666 // 100%-0% in 1/60th, i.e in 1sec

///|
pub struct Player {
  // manually doing "struct of arrays" optimization,
  // mostly as an experiment to see how tedious it is,
  // but also because MoonBit doesn't allow "#valtype"s in other "#valtype"s
  sprite_sheet : @util.SpriteSheet
  sprite_sheet_rev : @util.SpriteSheet
  peer : @firefly.Peer?
  mut pos : Vec2
  mut angle : Angle
  mut speed_factor : Float
}

///|
pub fn Player::new(
  sprite_sheet : @util.SpriteSheet,
  sprite_sheet_rev : @util.SpriteSheet,
  peer : @firefly.Peer?,
  pos? : Vec2 = Vec2::default(),
) -> Player {
  Player::{
    sprite_sheet,
    sprite_sheet_rev,
    peer,
    pos,
    angle: Angle::default(),
    speed_factor: 0,
  }
}

///|
pub fn Player::update(self : Player) -> Unit {
  self.sprite_sheet.update()
  self.sprite_sheet_rev.update()
  if self.peer is Some(peer) {
    self.update_input(peer)
  }
  let dir = Vec2::from_angle(self.angle).scale(MOVE_MAX_SPEED)
  self.pos = self.pos + dir.scale(MOVE_MAX_SPEED * self.speed_factor)
}

///|
pub fn Player::update_input(self : Player, peer : @firefly.Peer) -> Unit {
  let pad_maybe = @firefly.read_pad(peer~)
  guard pad_maybe is Some(pad) else {
    self.speed_factor = @util.move_towards(
      self.speed_factor,
      0,
      MOVE_DEACCELERATION,
    )
    return
  }
  let radius_percentage = pad.radius() / 1000.0
  // multiply by 1.2 and clamp so that there's some threshold to where
  // top speed is on the input
  let target_speed_factor = @cmp.minimum(radius_percentage * 1.2, 1.0)
  if target_speed_factor > self.speed_factor {
    self.speed_factor = @util.move_towards(
      self.speed_factor,
      target_speed_factor,
      MOVE_ACCELERATION,
    )
  } else {
    self.speed_factor = @util.move_towards(
      self.speed_factor,
      target_speed_factor,
      MOVE_DEACCELERATION,
    )
  }
  if pad.x != 0 && pad.y != 0 {
    let rotation_speed = @util.lerp(
      rotation_speed_when_still.to_rad(),
      rotation_speed_when_moving.to_rad(),
      self.speed_factor,
    )
    let new_angle = @util.rotate_toward(
      self.angle,
      pad.azimuth(),
      Angle::from_rad(rotation_speed),
    )
    self.angle = new_angle
  }
}

///|
pub fn Player::render(self : Player, world : World) -> Unit {
  let point = world.camera.world_vec2_to_camera_space(self.pos)
  // draw shadow
  @firefly.draw_circle(
    point + Point::new(-2, 2),
    5,
    @firefly.Style::new_solid(DarkGray),
  )
  // draw arrow of movement direction
  let dir = Vec2::from_angle(self.angle)
  @firefly.draw_line(
    point + dir.scale(6).to_point(),
    point + dir.scale(9).to_point(),
    @firefly.LineStyle::new(Black),
  )
  // draw sprite
  if Some(world.me) == self.peer {
    let is_looking_left = Angle::from_rad(@micromath.PI)
      .difference(self.angle)
      .to_rad()
      .abs() <
      @micromath.FRAC_PI_2
    if is_looking_left {
      self.sprite_sheet_rev.draw(point - Point::new(4, 5))
    } else {
      self.sprite_sheet.draw(point - Point::new(4, 5))
    }
  } else {
    @firefly.draw_circle(
      point - Point::new(2, 2),
      5,
      @firefly.Style::new_solid(Red),
    )
  }
}
