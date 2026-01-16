///|
const SCREEN_WIDTH_HALF_F : Float = 120

///|
const SCREEN_HEIGHT_HALF_F : Float = 80

///|
pub struct Camera {
  mut pos : Vec2
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
  let maybe_idx = world.player_system.players_peer.search(my_peer)
  guard maybe_idx is Some(idx) else { return }
  // Focus on player
  self.pos = world.player_system.players_pos[idx] -
    Vec2::new(SCREEN_WIDTH_HALF_F, SCREEN_HEIGHT_HALF_F)
}
