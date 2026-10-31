// Dataset constants and the settings of one generator run.

pub const SIZE: usize = 512;
pub const SHAPES_PER_IMAGE: usize = 20;
pub const R_MIN: i32 = 8; // scattered layout: shape "radius" range in pixels
pub const R_MAX: i32 = 30;
pub const SINGLE_R_MIN: i32 = 100; // single layout: shape "radius" range in pixels
pub const SINGLE_R_MAX: i32 = 240;
pub const EVAL_SEED: u64 = 20240601; // fixed, distinct from anything used for training data

pub struct Config {
    pub out_dir: String, // the images go to <out_dir>/train/<class>/ or <out_dir>/eval/<class>/
    pub images_per_class: usize,
    pub seed: u64,
}

impl Config {
    pub fn train_default(seed: u64) -> Config {
        Config { out_dir: "data_conv_std".to_string(), images_per_class: 1000, seed }
    }

    pub fn eval_default() -> Config {
        Config { out_dir: "data_conv_std".to_string(), images_per_class: 10, seed: EVAL_SEED }
    }
}
