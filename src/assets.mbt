///|
struct Assets {
  field : @firefly.Image
  firefly_sheet : FixedArray[@firefly.SubImage]
} derive(Default)

///|
fn Assets::load() -> Assets {
  Assets::{
    field: @firefly.load_file("field").unwrap().as_image(),
    firefly_sheet: Assets::load_firefly(),
  }
}

///|
fn Assets::load_firefly() -> FixedArray[@firefly.SubImage] {
  let img = @firefly.load_file("firefly").unwrap().as_image()
  let count = 7
  let size = Size::new(9, 10)
  let arr = FixedArray::make(count, @firefly.SubImage::default())
  for i in 0..<count {
    arr[i] = img.sub(Point::new(size.w * i, 0), size)
  }
  arr
}
