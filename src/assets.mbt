///|
struct Assets {
  field : @firefly.Image
  firefly_sheet : @util.SpriteSheet
}

///|
fn Assets::load() -> Assets {
  Assets::{
    field: @firefly.load_file("field").unwrap().as_image(),
    firefly_sheet: @util.SpriteSheet::split_image_by_count(
      @firefly.load_file("firefly").unwrap().as_image(),
      Size::new(7, 1),
    ),
  }
}
