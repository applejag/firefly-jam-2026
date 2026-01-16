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
    ..racebattle_world.val,
    player_system: @racebattle.PlayerSystem::new(@firefly.get_peers()),
  }
}

///|
pub fn update() -> Unit {
  racebattle_world.val.update()
}

///|
pub fn render() -> Unit {

}
