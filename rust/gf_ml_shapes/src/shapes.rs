// Shape classes and drawing of a single randomly colored/rotated shape.

use crate::canvas::{Canvas, Color, Point};
use crate::rng::Rng;

pub enum Class {
    Triangles,
    Circles,
    Squares,
    Lines,
}

impl Class {
    pub const ALL: [Class; 4] = [Class::Triangles, Class::Circles, Class::Squares, Class::Lines];

    /// Class name, also the name of its output directory.
    pub fn name(&self) -> &'static str {
        match self {
            Class::Triangles => "triangles",
            Class::Circles => "circles",
            Class::Squares => "squares",
            Class::Lines => "lines",
        }
    }
}

/// Random hue, high saturation/value so shapes stay visible on black.
pub fn random_color(rng: &mut Rng) -> Color {
    let hue = rng.float();
    let saturation = rng.uniform(0.6, 1.0);
    let value = rng.uniform(0.7, 1.0);
    hsv_to_rgb(hue, saturation, value)
}

fn hsv_to_rgb(h: f64, s: f64, v: f64) -> Color {
    let sector = (h * 6.0).floor();
    let f = h * 6.0 - sector;
    let p = v * (1.0 - s);
    let q = v * (1.0 - s * f);
    let t = v * (1.0 - s * (1.0 - f));
    let (r, g, b) = match sector as i32 % 6 {
        0 => (v, t, p),
        1 => (q, v, p),
        2 => (p, v, t),
        3 => (p, q, v),
        4 => (t, p, v),
        _ => (v, p, q),
    };
    [(r * 255.0) as u8, (g * 255.0) as u8, (b * 255.0) as u8]
}

/// The points rotated by `deg` degrees around (cx, cy).
pub fn rotated(points: &[Point], cx: f64, cy: f64, deg: f64) -> Vec<Point> {
    let t = deg.to_radians();
    let (s, c) = t.sin_cos();
    let mut out = Vec::with_capacity(points.len());
    for p in points {
        let dx = p.0 - cx;
        let dy = p.1 - cy;
        out.push((cx + dx * c - dy * s, cy + dx * s + dy * c));
    }
    out
}

/// Draws one shape of `class` centered at (cx, cy) with "radius" r.
/// `rotate_all`: true = every class gets a random rotation (training set);
/// false = only lines are rotated, squares and triangles stay upright (eval set).
pub fn draw_shape(canvas: &mut Canvas, rng: &mut Rng, class: &Class, cx: i32, cy: i32, r: i32, rotate_all: bool) {
    let color = random_color(rng);
    let angle = match class {
        Class::Lines => rng.uniform(0.0, 360.0),
        _ if rotate_all => rng.uniform(0.0, 360.0),
        _ => 0.0,
    };
    let (x, y, rf) = (cx as f64, cy as f64, r as f64);

    match class {
        Class::Circles => canvas.fill_circle(cx, cy, r, &color),
        Class::Squares => {
            let pts = [(x - rf, y - rf), (x + rf, y - rf), (x + rf, y + rf), (x - rf, y + rf)];
            canvas.fill_polygon(&rotated(&pts, x, y, angle), &color);
        }
        Class::Lines => {
            let w = ((r / 5).max(2)) as f64; // line thickness in pixels
            let pts = [(x - rf, y - w / 2.0), (x + rf, y - w / 2.0), (x + rf, y + w / 2.0), (x - rf, y + w / 2.0)];
            canvas.fill_polygon(&rotated(&pts, x, y, angle), &color);
        }
        Class::Triangles => {
            let pi = std::f64::consts::PI;
            let mut pts = Vec::with_capacity(3);
            for a in [-pi / 2.0, pi / 6.0, 5.0 * pi / 6.0] {
                pts.push((x + rf * a.cos(), y + rf * a.sin()));
            }
            canvas.fill_polygon(&rotated(&pts, x, y, angle), &color);
        }
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn hsv_primary_colors() {
        assert_eq!(hsv_to_rgb(0.0, 1.0, 1.0), [255, 0, 0]);
        assert_eq!(hsv_to_rgb(1.0 / 3.0, 1.0, 1.0)[1], 255);
    }

    #[test]
    fn every_class_draws_something() {
        for class in &Class::ALL {
            let mut canvas = Canvas::new(128, 128);
            let mut rng = Rng::new(7);
            draw_shape(&mut canvas, &mut rng, class, 64, 64, 30, true);
            assert!(canvas.pixels().iter().any(|&b| b != 0), "{} drew nothing", class.name());
        }
    }
}
