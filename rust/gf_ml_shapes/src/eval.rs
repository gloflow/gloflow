// Eval set (separate from the training data): SHAPES_PER_IMAGE random shapes per image, fixed seed.
// Squares and triangles stay upright, only lines are rotated.
//
// Output: <out_dir>/eval/<class>/<NN>.png  (512x512, RGB, colored shapes on black)

use std::fs;

use crate::canvas::Canvas;
use crate::config::*;
use crate::rng::Rng;
use crate::shapes::{draw_shape, Class};

pub fn generate_eval(config: &Config) -> Result<(), String> {
    let mut rng = Rng::new(config.seed);
    let size = SIZE as i32;

    for class in &Class::ALL {
        let class_dir = format!("{}/eval/{}", config.out_dir, class.name());
        fs::create_dir_all(&class_dir).map_err(|e| format!("create dir {}: {}", class_dir, e))?;

        for i in 0..config.images_per_class {
            let mut canvas = Canvas::new(SIZE, SIZE);
            for _ in 0..SHAPES_PER_IMAGE {
                let r = rng.int(R_MIN, R_MAX);
                let cx = rng.int(r, size - r);
                let cy = rng.int(r, size - r);
                draw_shape(&mut canvas, &mut rng, class, cx, cy, r, false);
            }
            canvas.save_png(&format!("{}/{:02}.png", class_dir, i))?;
        }
        println!("eval: {} done ({} images)", class.name(), config.images_per_class);
    }
    Ok(())
}
