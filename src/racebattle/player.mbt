///|
/// Rotation speed. Angle is the rotation speed per second
let rotation : Angle = Angle::from_deg(180.0 / 60.0)

///|
const MOVE_SPEED : Float = 5

///|
pub struct PlayerSystem {
  // manually doing "struct of arrays" optimization,
  // mostly as an experiment to see how tedious it is,
  // but also because MoonBit doesn't allow "#valtype"s in other "#valtype"s
  players_peer : ReadOnlyArray[@firefly.Peer]
  players_pos : FixedArray[Vec2]
  players_angle : FixedArray[Angle]
} derive(Default)

///|
pub fn PlayerSystem::new(peers : @firefly.Peers) -> PlayerSystem {
  let peers = peers.to_array()
  PlayerSystem::{
    players_peer: ReadOnlyArray::from_array(peers),
    players_pos: FixedArray::make(peers.length(), Vec2::default()),
    players_angle: FixedArray::make(peers.length(), Angle::default()),
  }
}

///|
pub fn PlayerSystem::update(self : PlayerSystem) -> Unit {
  for idx, peer in self.players_peer {
    let pad_maybe = @firefly.read_pad(peer~)
    guard pad_maybe is Some(pad) else { continue }
    guard pad.x != 0 && pad.y != 0 else { continue }
    let pos = self.players_pos[idx]
    let angle = self.players_angle[idx]
    let new_angle = @util.rotate_toward(angle, pad.azimuth(), rotation)
    let dir = Vec2::from_angle(new_angle).scale(MOVE_SPEED)
    self.players_pos[idx] = pos + dir
    self.players_angle[idx] = new_angle
  }
}
