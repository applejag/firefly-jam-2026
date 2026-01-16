///|
pub struct SpriteSheet {
  sprites : FixedArray[@firefly.SubImage]
  mut index : Int
  mut time : Int
}

///|
pub fn SpriteSheet::split_image_by_count(
  image : @firefly.Image,
  counts : Size,
  image_size? : Size = image.size(),
) -> SpriteSheet {
  let sprite_size = Size::new(image_size.w / counts.w, image_size.h / counts.h)
  SpriteSheet::split_image_by_size(image, sprite_size, image_size~)
}

///|
pub fn SpriteSheet::split_image_by_size(
  image : @firefly.Image,
  sprite_size : Size,
  image_size? : Size = image.size(),
) -> SpriteSheet {
  let count_per_axis = Size::new(
    image_size.w / sprite_size.w,
    image_size.h / sprite_size.h,
  )
  let count = count_per_axis.w * count_per_axis.h
  let arr = FixedArray::make(count, @firefly.SubImage::default())
  for i in 0..<count {
    arr[i] = image.sub(
      Point::new(
        i % count_per_axis.w * sprite_size.w,
        i / count_per_axis.h * sprite_size.h,
      ),
      sprite_size,
    )
  }
  SpriteSheet::{ sprites: arr, index: 0, time: 0 }
}

///|
/// Updates a sprite sheet's animation cycling.
/// If `fps` is 0 or negative, then the animation will never change.
pub fn SpriteSheet::update(self : SpriteSheet, fps : Int) -> Unit {
  guard fps > 0 else { return }
  self.time += 1
  //60 / fps
  //20
  self.index = (self.index + 1) % self.sprites.length()
}

///|
#inline
pub fn SpriteSheet::draw(self : SpriteSheet, point : Point) -> Unit {
  self.sprites[self.index].draw(point)
}
