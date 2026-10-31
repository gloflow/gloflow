// RGB canvas on top of tiny-skia (pure Rust 2D rasterizer): filled circles and polygons, PNG saving.

use std::fs::File;
use std::io::BufWriter;

use tiny_skia::{FillRule, Paint, PathBuilder, Pixmap, Transform};

pub type Point = (f64, f64);
pub type Color = [u8; 3];

// Anti-aliased edges. false = hard (aliased) edges.
const ANTI_ALIAS: bool = true;

pub struct Canvas {
    width: usize,
    height: usize,
    pixmap: Pixmap, // RGBA premultiplied; always fully opaque here, so RGB == the stored color
}

impl Canvas {
    /// New black canvas.
    pub fn new(width: usize, height: usize) -> Canvas {
        let mut pixmap = Pixmap::new(width as u32, height as u32).expect("canvas size must be non-zero");
        pixmap.fill(tiny_skia::Color::BLACK);
        Canvas { width, height, pixmap }
    }

    /// RGB bytes, row by row: 3 bytes per pixel.
    pub fn pixels(&self) -> Vec<u8> {
        let mut out = Vec::with_capacity(self.width * self.height * 3);
        for p in self.pixmap.data().chunks_exact(4) {
            out.extend_from_slice(&p[..3]);
        }
        out
    }

    pub fn get_pixel(&self, x: usize, y: usize) -> Color {
        let i = (y * self.width + x) * 4;
        let d = self.pixmap.data();
        [d[i], d[i + 1], d[i + 2]]
    }

    fn paint(color: &Color) -> Paint<'static> {
        let mut paint = Paint::default();
        paint.set_color_rgba8(color[0], color[1], color[2], 255);
        paint.anti_alias = ANTI_ALIAS;
        paint
    }

    /// Filled circle of radius r centered at the point (cx, cy) (continuous coordinates, pixel
    /// (x, y) covers x..x+1, y..y+1). Parts outside of the canvas are clipped.
    pub fn fill_circle(&mut self, cx: i32, cy: i32, r: i32, color: &Color) {
        let path = PathBuilder::from_circle(cx as f32, cy as f32, r as f32);
        if let Some(path) = path {
            self.pixmap.fill_path(&path, &Self::paint(color), FillRule::Winding, Transform::identity(), None);
        }
    }

    /// Filled polygon with the same continuous coordinates as `fill_circle`.
    /// Parts outside of the canvas are clipped.
    pub fn fill_polygon(&mut self, points: &[Point], color: &Color) {
        if points.len() < 3 {
            return;
        }
        let mut pb = PathBuilder::new();
        for (i, p) in points.iter().enumerate() {
            let (x, y) = (p.0 as f32, p.1 as f32);
            if i == 0 {
                pb.move_to(x, y);
            } else {
                pb.line_to(x, y);
            }
        }
        pb.close();
        if let Some(path) = pb.finish() {
            self.pixmap.fill_path(&path, &Self::paint(color), FillRule::Winding, Transform::identity(), None);
        }
    }

    /// Writes the canvas as an 8-bit RGB PNG file.
    pub fn save_png(&self, path: &str) -> Result<(), String> {
        let file = File::create(path).map_err(|e| format!("create {}: {}", path, e))?;
        let mut encoder = png::Encoder::new(BufWriter::new(file), self.width as u32, self.height as u32);
        encoder.set_color(png::ColorType::Rgb);
        encoder.set_depth(png::BitDepth::Eight);
        let mut writer = encoder.write_header().map_err(|e| format!("png header {}: {}", path, e))?;
        writer.write_image_data(&self.pixels()).map_err(|e| format!("png data {}: {}", path, e))?;
        Ok(())
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    const WHITE: Color = [255, 255, 255];

    fn count_lit(canvas: &Canvas) -> usize {
        canvas.pixels().chunks(3).filter(|p| p[0] != 0).count()
    }

    #[test]
    fn circle_is_filled_around_its_center() {
        let mut canvas = Canvas::new(64, 64);
        canvas.fill_circle(32, 32, 10, &WHITE);
        assert_eq!(canvas.get_pixel(32, 32), WHITE);
        assert_eq!(canvas.get_pixel(32, 39), WHITE);
        assert_eq!(canvas.get_pixel(32, 44), [0, 0, 0]);
        // lit pixels (edge pixels included) are close to the circle area
        let n = count_lit(&canvas) as f64;
        let ideal = std::f64::consts::PI * 10.0 * 10.0;
        assert!((n - ideal).abs() < ideal * 0.15, "{} lit, ideal {}", n, ideal);
    }

    #[test]
    fn axis_aligned_square_covers_exact_pixels() {
        let mut canvas = Canvas::new(64, 64);
        let pts = [(10.0, 10.0), (20.0, 10.0), (20.0, 20.0), (10.0, 20.0)];
        canvas.fill_polygon(&pts, &WHITE);
        assert_eq!(count_lit(&canvas), 10 * 10);
        assert_eq!(canvas.get_pixel(10, 10), WHITE);
        assert_eq!(canvas.get_pixel(20, 20), [0, 0, 0]);
    }

    #[test]
    fn shapes_are_clipped_at_the_canvas_border() {
        let mut canvas = Canvas::new(16, 16);
        canvas.fill_circle(0, 0, 5, &WHITE);
        let pts = [(10.0, 10.0), (30.0, 10.0), (30.0, 30.0)];
        canvas.fill_polygon(&pts, &WHITE);
        assert!(count_lit(&canvas) > 0);
    }
}
