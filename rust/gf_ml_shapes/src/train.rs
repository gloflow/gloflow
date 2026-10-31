// Training set: images_per_class images per class, split into layouts.
//
// Layouts (percentages must sum to 100):
//   single     one large instance of the shape (teaches what a single shape looks like)
//   scattered  SHAPES_PER_IMAGE small randomly placed, overlapping shapes
//
// Output: <out_dir>/train/<class>/<layout>_<NNNN>.png  (512x512, RGB, colored shapes on black)
//         <out_dir>/train/<class>/manifest.json         (layout names, their percentage and image count)

use std::fs;

use crate::canvas::Canvas;
use crate::config::*;
use crate::rng::Rng;
use crate::shapes::{draw_shape, Class};

struct Layout {
    name: &'static str,
    percent: usize,
    description: String,
}

fn layouts() -> Vec<Layout> {
    vec![
        Layout {
            name: "single",
            percent: 50,
            description: "one large instance of the shape, random color/rotation/position".to_string(),
        },
        Layout {
            name: "scattered",
            percent: 50,
            description: format!("{} small randomly placed, possibly overlapping instances of the shape", SHAPES_PER_IMAGE),
        },
    ]
}

/// Image count per layout; the last layout gets the rounding remainder.
fn layout_counts(layouts: &[Layout], images_per_class: usize) -> Vec<usize> {
    let percent_sum: usize = layouts.iter().map(|l| l.percent).sum();
    assert_eq!(percent_sum, 100, "layout percentages must sum to 100");
    let mut counts: Vec<usize> = layouts.iter().map(|l| images_per_class * l.percent / 100).collect();
    let assigned: usize = counts.iter().sum();
    if let Some(last) = counts.last_mut() {
        *last += images_per_class - assigned;
    }
    counts
}

// random position so that the whole shape of radius r stays inside the image
fn random_center(rng: &mut Rng, r: i32) -> (i32, i32) {
    let size = SIZE as i32;
    (rng.int(r, size - r), rng.int(r, size - r))
}

fn draw_layout(canvas: &mut Canvas, rng: &mut Rng, class: &Class, layout_name: &str) {
    match layout_name {
        "single" => {
            let r = rng.int(SINGLE_R_MIN, SINGLE_R_MAX);
            let (cx, cy) = random_center(rng, r);
            draw_shape(canvas, rng, class, cx, cy, r, true);
        }
        _ => {
            // "scattered"
            for _ in 0..SHAPES_PER_IMAGE {
                let r = rng.int(R_MIN, R_MAX);
                let (cx, cy) = random_center(rng, r);
                draw_shape(canvas, rng, class, cx, cy, r, true);
            }
        }
    }
}

fn manifest_json(class: &Class, images_per_class: usize, layouts: &[Layout], counts: &[usize]) -> String {
    let mut items = Vec::new();
    for (layout, count) in layouts.iter().zip(counts) {
        items.push(format!(
            "    {{\n      \"name\": \"{}\",\n      \"percent\": {},\n      \"images\": {},\n      \"file_prefix\": \"{}_\",\n      \"description\": \"{}\"\n    }}",
            layout.name, layout.percent, count, layout.name, layout.description
        ));
    }
    format!(
        "{{\n  \"class\": \"{}\",\n  \"total_images\": {},\n  \"layouts\": [\n{}\n  ]\n}}\n",
        class.name(),
        images_per_class,
        items.join(",\n")
    )
}

// deletes stale .png files of earlier runs
fn remove_old_images(dir: &str) -> Result<(), String> {
    let entries = fs::read_dir(dir).map_err(|e| format!("read dir {}: {}", dir, e))?;
    for entry in entries {
        let path = entry.map_err(|e| e.to_string())?.path();
        if path.extension().map_or(false, |ext| ext == "png") {
            fs::remove_file(&path).map_err(|e| format!("remove {}: {}", path.display(), e))?;
        }
    }
    Ok(())
}

pub fn generate_train(config: &Config) -> Result<(), String> {
    let mut rng = Rng::new(config.seed);
    let layouts = layouts();
    let counts = layout_counts(&layouts, config.images_per_class);

    for class in &Class::ALL {
        let class_dir = format!("{}/train/{}", config.out_dir, class.name());
        fs::create_dir_all(&class_dir).map_err(|e| format!("create dir {}: {}", class_dir, e))?;
        remove_old_images(&class_dir)?;

        for (layout, count) in layouts.iter().zip(&counts) {
            for i in 0..*count {
                let mut canvas = Canvas::new(SIZE, SIZE);
                draw_layout(&mut canvas, &mut rng, class, layout.name);
                canvas.save_png(&format!("{}/{}_{:04}.png", class_dir, layout.name, i))?;
            }
        }

        let manifest_path = format!("{}/manifest.json", class_dir);
        let json = manifest_json(class, config.images_per_class, &layouts, &counts);
        fs::write(&manifest_path, json).map_err(|e| format!("write {}: {}", manifest_path, e))?;
        println!("train: {} done ({} images)", class.name(), config.images_per_class);
    }
    Ok(())
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn counts_sum_to_total_with_remainder_in_last() {
        let l = layouts();
        assert_eq!(layout_counts(&l, 1000), vec![500, 500]);
        assert_eq!(layout_counts(&l, 7), vec![3, 4]);
    }

    #[test]
    fn manifest_lists_layouts() {
        let l = layouts();
        let json = manifest_json(&Class::Circles, 10, &l, &[5, 5]);
        assert!(json.contains("\"class\": \"circles\""));
        assert!(json.contains("\"file_prefix\": \"scattered_\""));
    }
}
