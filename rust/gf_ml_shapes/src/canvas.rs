// RGB pixel buffer with the two fill operations the shapes need, and PNG saving.

use std::fs::File;
use std::io::BufWriter;

pub type Point = (f64, f64);
pub type Color = [u8; 3];

pub struct Canvas {
    width: usize,
    height: usize,
    pixels: Vec<u8>, // RGB, row by row: 3 bytes per pixel
}

impl Canvas {
    /// New black canvas.
    pub fn new(width: usize, height: usize) -> Canvas {
        Canvas { width, height, pixels: vec![0; width * height * 3] }
    }

    pub fn pixels(&self) -> &[u8] {
        &self.pixels
    }

    pub fn get_pixel(&self, x: usize, y: usize) -> &[u8] {
        let i = (y * self.width + x) * 3;
        &self.pixels[i..i + 3]
    }

    // pixels outside of the canvas are ignored
    fn set_pixel(&mut self, x: i32, y: i32, color: &Color) {
        if x < 0 || y < 0 || x >= self.width as i32 || y >= self.height as i32 {
            return;
        }
        let i = (y as usize * self.width + x as usize) * 3;
        self.pixels[i..i + 3].copy_from_slice(color);
    }

    /// Filled circle. Covers the pixels cx-r ..= cx+r (same as a Pillow ellipse box).
    pub fn fill_circle(&mut self, cx: i32, cy: i32, r: i32, color: &Color) {
        let limit = (r as f64 + 0.5) * (r as f64 + 0.5);
        for y in (cy - r)..=(cy + r) {
            for x in (cx - r)..=(cx + r) {
                let dx = (x - cx) as f64;
                let dy = (y - cy) as f64;
                if dx * dx + dy * dy <= limit {
                    self.set_pixel(x, y, color);
                }
            }
        }
    }

    /// Filled CONVEX polygon (all shapes here are convex). A point (x, y) is inside
    /// when the pixel at integer coordinates (x, y) lies within the polygon edges.
    pub fn fill_polygon(&mut self, points: &[Point], color: &Color) {
        if points.len() < 3 {
            return;
        }
        let mut y_min = f64::MAX;
        let mut y_max = f64::MIN;
        for p in points {
            y_min = y_min.min(p.1);
            y_max = y_max.max(p.1);
        }
        for y in (y_min.ceil() as i32)..=(y_max.floor() as i32) {
            if y < 0 || y >= self.height as i32 {
                continue;
            }
            // convex: the row is filled from the left-most to the right-most edge crossing
            let mut x_min = f64::MAX;
            let mut x_max = f64::MIN;
            for i in 0..points.len() {
                let a = &points[i];
                let b = &points[(i + 1) % points.len()];
                let yf = y as f64;
                if yf < a.1.min(b.1) || yf > a.1.max(b.1) {
                    continue; // edge does not cross this row
                }
                if a.1 == b.1 {
                    // horizontal edge: both ends count
                    x_min = x_min.min(a.0.min(b.0));
                    x_max = x_max.max(a.0.max(b.0));
                } else {
                    let t = (yf - a.1) / (b.1 - a.1);
                    let x = a.0 + t * (b.0 - a.0);
                    x_min = x_min.min(x);
                    x_max = x_max.max(x);
                }
            }
            if x_min > x_max {
                continue;
            }
            for x in (x_min.ceil() as i32)..=(x_max.floor() as i32) {
                self.set_pixel(x, y, color);
            }
        }
    }

    /// Writes the canvas as an 8-bit RGB PNG file.
    pub fn save_png(&self, path: &str) -> Result<(), String> {
        let file = File::create(path).map_err(|e| format!("create {}: {}", path, e))?;
        let mut encoder = png::Encoder::new(BufWriter::new(file), self.width as u32, self.height as u32);
        encoder.set_color(png::ColorType::Rgb);
        encoder.set_depth(png::BitDepth::Eight);
        let mut writer = encoder.write_header().map_err(|e| format!("png header {}: {}", path, e))?;
        writer.write_image_data(&self.pixels).map_err(|e| format!("png data {}: {}", path, e))?;
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
        assert_eq!(canvas.get_pixel(32, 32), &WHITE);
        assert_eq!(canvas.get_pixel(32, 42), &WHITE);
        assert_eq!(canvas.get_pixel(32, 44), &[0, 0, 0]);
        let n = count_lit(&canvas) as f64;
        let ideal = std::f64::consts::PI * 10.5 * 10.5;
        assert!((n - ideal).abs() < ideal * 0.05);
    }

    #[test]
    fn axis_aligned_square_covers_inclusive_pixels() {
        let mut canvas = Canvas::new(64, 64);
        let pts = [(10.0, 10.0), (20.0, 10.0), (20.0, 20.0), (10.0, 20.0)];
        canvas.fill_polygon(&pts, &WHITE);
        assert_eq!(count_lit(&canvas), 11 * 11);
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
