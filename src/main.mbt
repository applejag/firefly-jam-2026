///|
using @firefly {type Size}

///|
using @util {type Vec2}

///|
let racebattle_world : Ref[@racebattle.World] = Ref::new(
  @racebattle.World::default(),
)

///|
fn main {

}

///|
pub fn boot() -> Unit {
  let assets = Assets::load()
  racebattle_world.val = @racebattle.World::{
    field_img: assets.field,
    players: [
      @racebattle.Player::new(
        assets.firefly_sheet,
        Some(@firefly.get_me()),
        pos=Vec2::new(80, 40),
      ),
      //@racebattle.Player::new(assets.firefly_sheet, None),
      //@racebattle.Player::new(assets.firefly_sheet, None),
    ],
    camera: @racebattle.Camera::default(),
    me: @firefly.get_me(),
  }
}

///|
pub fn update() -> Unit {
  racebattle_world.val.update()
}

///|
pub fn render() -> Unit {
  racebattle_world.val.render()
}
