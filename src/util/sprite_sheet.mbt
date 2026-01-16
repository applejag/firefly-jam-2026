///|
pub struct SpriteSheet {
  sprites : FixedArray[@firefly.SubImage]
  mut index : Int
  mut time : Int
  ticks_per_frame : Int
}

///|
pub fn SpriteSheet::split_image_by_count(
  image : @firefly.Image,
  counts : Size,
  image_size? : Size = image.size(),
  fps? : Int = 30,
) -> SpriteSheet {
  let sprite_size = Size::new(image_size.w / counts.w, image_size.h / counts.h)
  SpriteSheet::split_image_by_size(image, sprite_size, image_size~, fps~)
}

///|
pub fn SpriteSheet::split_image_by_size(
  image : @firefly.Image,
  sprite_size : Size,
  image_size? : Size = image.size(),
  fps? : Int = 30,
) -> SpriteSheet {
  let count_per_axis = Size::new(
    image_size.w / sprite_size.w,
    image_size.h / sprite_size.h,
  )
  let count = count_per_axis.w * count_per_axis.h
  let arr = FixedArray::make(count, @firefly.SubImage::default())
  for i in 0..<count {
    let point = Point::new(
      i % count_per_axis.w * sprite_size.w,
      i / count_per_axis.w * sprite_size.h,
    )
    arr[i] = image.sub(point, sprite_size)
  }
  SpriteSheet::{ sprites: arr, index: 0, time: 0, ticks_per_frame: 60 / fps }
}

///|
/// Updates a sprite sheet's animation cycling.
/// If `fps` is 0 or negative, then the animation will never change.
pub fn SpriteSheet::update(self : SpriteSheet) -> Unit {
  self.time += 1
  if self.time >= self.ticks_per_frame {
    self.index = (self.index + 1) % self.sprites.length()
    self.time = 0
  }
}

///|
#inline
pub fn SpriteSheet::draw(self : SpriteSheet, point : Point) -> Unit {
  self.sprites[self.index].draw(point)
}
