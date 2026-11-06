//! gf_lang_render - renders a gf_lang scene (JSON) to a PNG, headless.
//!
//!   gf_lang_render [--scene FILE|-] --out out.png [--width 1024] [--height 768] [--samples 4]
//!
//! Scene JSON is read from FILE, or stdin if `--scene` is omitted or `-`.

use gf_lang_render::{encode_png, render, Options, Scene};
use std::io::Read;

fn usage() -> ! {
    eprintln!("usage: gf_lang_render [--scene FILE|-] --out out.png [--width N] [--height N] [--samples 1|2|4] [--fit]");
    std::process::exit(2);
}

fn run() -> Result<(), String> {
    let mut scene_path = String::from("-");
    let mut out_path: Option<String> = None;
    let mut opts = Options::default();

    let mut args = std::env::args().skip(1);
    while let Some(a) = args.next() {
        let mut val = |name: &str| args.next().unwrap_or_else(|| {
            eprintln!("{name} needs a value");
            usage()
        });
        match a.as_str() {
            "--scene" => scene_path = val("--scene"),
            "--out" => out_path = Some(val("--out")),
            "--width" => opts.width = val("--width").parse().map_err(|_| "bad --width")?,
            "--height" => opts.height = val("--height").parse().map_err(|_| "bad --height")?,
            "--fit" => opts.fit = true,
            "--samples" => opts.samples = val("--samples").parse().map_err(|_| "bad --samples")?,
            "-h" | "--help" => usage(),
            other => {
                eprintln!("unknown argument: {other}");
                usage()
            }
        }
    }
    let out_path = out_path.unwrap_or_else(|| usage());

    let mut bytes = Vec::new();
    if scene_path == "-" {
        std::io::stdin().read_to_end(&mut bytes).map_err(|e| e.to_string())?;
    } else {
        bytes = std::fs::read(&scene_path).map_err(|e| format!("{scene_path}: {e}"))?;
    }

    let scene = Scene::from_json(&bytes)?;
    if scene.truncated {
        eprintln!("warning: scene was truncated by the entity cap");
    }
    let image = render(&scene, &opts)?;
    let png = encode_png(&image)?;
    std::fs::write(&out_path, png).map_err(|e| format!("{out_path}: {e}"))?;
    eprintln!("wrote {out_path} ({}x{}, {} entities)", image.width, image.height, scene.entities.len());
    Ok(())
}

fn main() {
    if let Err(e) = run() {
        eprintln!("error: {e}");
        std::process::exit(1);
    }
}
