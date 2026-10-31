// Usage: gf_ml_shapes <train|eval> [--out DIR] [--images-per-class N] [--seed N]
//   train  writes DIR/train/<class>/  (default: 1000 images per class, random seed)
//   eval   writes DIR/eval/<class>/   (default: 10 images per class, fixed seed)
//   --out  default "data_conv_std"

use std::time::{SystemTime, UNIX_EPOCH};

use gf_ml_shapes::config::Config;
use gf_ml_shapes::eval::generate_eval;
use gf_ml_shapes::train::generate_train;

const USAGE: &str = "usage: gf_ml_shapes <train|eval> [--out DIR] [--images-per-class N] [--seed N]";

fn time_seed() -> u64 {
    SystemTime::now().duration_since(UNIX_EPOCH).map(|d| d.as_nanos() as u64).unwrap_or(1)
}

// applies the --flag value pairs in args onto config
fn apply_flags(config: &mut Config, args: &[String]) -> Result<(), String> {
    let mut i = 0;
    while i < args.len() {
        let value = args.get(i + 1).ok_or(format!("{} needs a value", args[i]))?;
        match args[i].as_str() {
            "--out" => config.out_dir = value.to_string(),
            "--images-per-class" => config.images_per_class = value.parse().map_err(|_| format!("bad number: {}", value))?,
            "--seed" => config.seed = value.parse().map_err(|_| format!("bad number: {}", value))?,
            other => return Err(format!("unknown option {}\n{}", other, USAGE)),
        }
        i += 2;
    }
    Ok(())
}

fn run(args: &[String]) -> Result<(), String> {
    let command = args.get(1).ok_or(USAGE.to_string())?;
    let flags = &args[2..];
    match command.as_str() {
        "train" => {
            let mut config = Config::train_default(time_seed());
            apply_flags(&mut config, flags)?;
            generate_train(&config)
        }
        "eval" => {
            let mut config = Config::eval_default();
            apply_flags(&mut config, flags)?;
            generate_eval(&config)
        }
        _ => Err(USAGE.to_string()),
    }
}

fn main() {
    let args: Vec<String> = std::env::args().collect();
    if let Err(msg) = run(&args) {
        eprintln!("ERROR!! - {}", msg);
        std::process::exit(1);
    }
}
