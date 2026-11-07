//! gf_lang_render - renders a gf_lang scene (JSON) to a PNG, headless.
//!
//!   gf_lang_render [--scene FILE|-] --out out.png [--width 1024] [--height 768] [--samples 4]
//!                  [--fit | --camera-pos x,y,z --camera-target x,y,z [--camera-up x,y,z] [--fov DEG]]
//!
//! Scene JSON is read from FILE, or stdin if `--scene` is omitted or `-`.
//! The camera flags render from an explicit view (ex. the dashboard's view); `--fit` frames all entities.
//! Without either, the scene camera is used.

use gf_lang_render::{encode_png, render, CameraView, Options, Scene};
use std::io::Read;

fn usage() -> ! {
    eprintln!(
        "usage: gf_lang_render [--scene FILE|-] --out out.png [--width N] [--height N] [--samples 1|2|4]\n\
         \x20                     [--fit | --camera-pos x,y,z --camera-target x,y,z [--camera-up x,y,z] [--fov DEG]]"
    );
    std::process::exit(2);
}

fn parse_vec3(name: &str, s: &str) -> Result<[f32; 3], String> {
    let parts: Vec<f32> = s
        .split(',')
        .map(|p| p.trim().parse::<f32>())
        .collect::<Result<_, _>>()
        .map_err(|_| format!("{name}: expected x,y,z numbers, got '{s}'"))?;
    <[f32; 3]>::try_from(parts).map_err(|_| format!("{name}: expected exactly 3 numbers, got '{s}'"))
}

fn run() -> Result<(), String> {
    let mut scene_path = String::from("-");
    let mut out_path: Option<String> = None;
    let mut opts = Options::default();

    let mut camera_pos: Option<[f32; 3]> = None;
    let mut camera_target: Option<[f32; 3]> = None;
    let mut camera_up: Option<[f32; 3]> = None;
    let mut camera_fov: Option<f32> = None;

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
            "--camera-pos" => camera_pos = Some(parse_vec3("--camera-pos", &val("--camera-pos"))?),
            "--camera-target" => camera_target = Some(parse_vec3("--camera-target", &val("--camera-target"))?),
            "--camera-up" => camera_up = Some(parse_vec3("--camera-up", &val("--camera-up"))?),
            "--fov" => camera_fov = Some(val("--fov").parse().map_err(|_| "bad --fov")?),
            "-h" | "--help" => usage(),
            other => {
                eprintln!("unknown argument: {other}");
                usage()
            }
        }
    }
    let out_path = out_path.unwrap_or_else(|| usage());

    // an explicit camera needs a position and a target; up defaults to +Y and fov to 75 (the browser engine's values)
    let any_camera_flag = camera_pos.is_some() || camera_target.is_some() || camera_up.is_some() || camera_fov.is_some();
    if any_camera_flag {
        match (camera_pos, camera_target) {
            (Some(position), Some(target)) => {
                opts.camera = Some(CameraView {
                    position,
                    target,
                    up: camera_up.unwrap_or([0.0, 1.0, 0.0]),
                    fov_deg: camera_fov.unwrap_or(75.0),
                });
            }
            _ => return Err("--camera-pos and --camera-target must be given together".into()),
        }
    }

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
