///|
const SCREEN_WIDTH_HALF_F : Float = 120

///|
const SCREEN_HEIGHT_HALF_F : Float = 80

///|
pub struct Camera {
  pos : Vec2
} derive(Default)

///|
#inline
fn Camera::world_vec2_to_camera_space(self : Camera, position : Vec2) -> Point {
  (position - self.pos).round().to_point()
}

///|
#inline
fn Camera::world_point_to_camera_space(self : Camera, point : Point) -> Point {
  point - self.pos.round().to_point()
}

///|
fn Camera::update(self : Camera, world : World) -> Unit {
  let my_peer = @firefly.get_me()
  let maybe_player = world.players.iter().find_first(p => p.peer == my_peer)
  if maybe_player is Some(player) {
    player.update()
  }
}
