///|
pub(all) struct World {
  field_img : @firefly.Image
  players : FixedArray[Player]
  camera : Camera
  me : @firefly.Peer
} derive(Default)

///|
pub fn World::update(self : World) -> Unit {
  for player in self.players {
    player.update()
  }
  self.camera.update(self)
}

///|
pub fn World::render(self : World) -> Unit {
  // background
  @firefly.draw_rect(
    Point::new(0, 0),
    @firefly.Size::new_screen(),
    @firefly.Style::new_solid(DarkGray),
  )
  self.field_img.draw(self.camera.world_point_to_camera_space(Point::new(0, 0)))
  // players
  let mut me : Player? = None
  for player in self.players {
    if player.peer == Some(self.me) {
      me = Some(player)
      continue
    }
    player.render(self)
  }
  // always render my player last (on top)
  if me is Some(player) {
    player.render(self)
  }
}
