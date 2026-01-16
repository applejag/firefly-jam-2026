///|
let racebattle_world : Ref[@racebattle.World] = Ref::new(
  @racebattle.World::default(),
)

///|
fn main {

}

///|
pub fn boot() -> Unit {
  //let assets = Assets::load()
  racebattle_world.val = @racebattle.World::{
    ..racebattle_world.val,
    players: [@racebattle.Player::new(@firefly.get_me())],
  }
}

///|
pub fn update() -> Unit {
  racebattle_world.val.update()
}

///|
pub fn render() -> Unit {

}
