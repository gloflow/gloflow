// Synthetic shapes dataset: randomly colored shapes on a black background, saved as PNG files.
//
//   rng     - tiny seedable random number generator
//   canvas  - RGB pixel buffer + filling of circles/polygons + PNG saving
//   shapes  - shape classes and drawing of one shape
//   config  - sizes, radii and the run settings
//   train   - training set (layouts: "single", "scattered") + manifest.json per class
//   eval    - eval set (fixed seed)

pub mod canvas;
pub mod config;
pub mod eval;
pub mod rng;
pub mod shapes;
pub mod train;
