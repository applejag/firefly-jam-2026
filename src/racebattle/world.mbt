///|
pub(all) struct World {
  player_system : PlayerSystem
  camera : Camera
} derive(Default)

///|
pub fn World::update(self : World) -> Unit {
  self.player_system.update()
  self.camera.update(self)
}
