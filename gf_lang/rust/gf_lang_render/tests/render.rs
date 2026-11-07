//! Tests that don't need a GPU (scene parsing, shader translation) and a render test that does.
//! The render test is skipped if no adapter is available (set GF_RENDER_REQUIRE_GPU=1 to make it fail instead).

use gf_lang_render::scene::{Material, Scene, ShaderDef};
use gf_lang_render::{camera_setup, render, shader, CameraView, Options};
use std::collections::HashMap;

const SCENE: &str = r#"{
  "version": 1,
  "background": [1, 1, 1],
  "camera": {"position": [0, 0, 5], "target": [0, 0, 0], "fov_deg": 60, "near": 0.1, "far": 100},
  "shaders": {},
  "entities": [
    {"kind": "cube", "position": [0,0,0], "rotation": [0,0,0], "scale": [1,1,1],
     "color": [1, 0, 0], "wireframe": false, "material": {"type": "basic"}}
  ],
  "truncated": false
}"#;

#[test]
fn parses_scene() {
    let s = Scene::from_json(SCENE.as_bytes()).unwrap();
    assert_eq!(s.entities.len(), 1);
    assert!(matches!(s.entities[0].material, Material::Basic));
}

#[test]
fn rejects_unknown_version() {
    let bad = SCENE.replace("\"version\": 1", "\"version\": 99");
    assert!(Scene::from_json(bad.as_bytes()).is_err());
}

fn shader_def() -> ShaderDef {
    serde_json::from_str(
        r#"{
        "uniforms": [{"name": "i", "type": "float", "default": 0.1}, {"name": "tint", "type": "vec3", "default": [1,0,0]}],
        "vertex": "varying vec3 n;\nvarying vec3 p;\nvoid main() { n = normal; p = position; gl_Position = projectionMatrix * modelViewMatrix * vec4(position, 1.0); }",
        "fragment": "precision highp float;\nuniform float i; // comment\nvarying vec3 p;\nvarying vec3 n;\nvoid main() { gl_FragColor = vec4(p * i, 1.0); }"
    }"#,
    )
    .unwrap()
}

#[test]
fn translates_three_style_glsl() {
    let t = shader::translate(&shader_def()).unwrap();
    assert!(t.vertex.starts_with("#version 450"));
    assert!(t.vertex.contains("layout(location = 0) out vec3 n;"));
    assert!(t.vertex.contains("layout(location = 1) out vec3 p;"));
    // fragment varyings take their location from the vertex shader, not their own order
    assert!(t.fragment.contains("layout(location = 1) in vec3 p;"));
    assert!(t.fragment.contains("layout(location = 0) in vec3 n;"));
    assert!(t.fragment.contains("gf_FragColor"));
    assert!(!t.fragment.contains("precision"));
    // std140: float at 0, vec3 at 16
    assert_eq!(t.uniforms[0].offset, 0);
    assert_eq!(t.uniforms[1].offset, 16);
    assert_eq!(t.block_size, 32);
}

#[test]
fn packs_uniforms() {
    let t = shader::translate(&shader_def()).unwrap();
    let mut vals = HashMap::new();
    vals.insert("i".to_string(), serde_json::json!(0.5));
    let buf = shader::pack_uniforms(&t, &vals).unwrap();
    assert_eq!(f32::from_le_bytes(buf[0..4].try_into().unwrap()), 0.5);
    // default for tint
    assert_eq!(f32::from_le_bytes(buf[16..20].try_into().unwrap()), 1.0);
}

#[test]
fn rejects_unmatched_varying() {
    let mut def = shader_def();
    def.fragment = "varying vec3 nope;\nvoid main() { gl_FragColor = vec4(nope, 1.0); }".into();
    assert!(shader::translate(&def).is_err());
}

#[test]
fn renders_cube() {
    let scene = Scene::from_json(SCENE.as_bytes()).unwrap();
    let opts = Options { width: 64, height: 64, samples: 1, fit: false, camera: None };
    let img = match render(&scene, &opts) {
        Ok(i) => i,
        Err(e) if e.contains("no GPU adapter") && std::env::var("GF_RENDER_REQUIRE_GPU").is_err() => {
            eprintln!("skipping: {e}");
            return;
        }
        Err(e) => panic!("{e}"),
    };
    assert_eq!(img.rgba.len(), 64 * 64 * 4);
    let px = |x: usize, y: usize| &img.rgba[(y * 64 + x) * 4..(y * 64 + x) * 4 + 4];
    assert_eq!(px(32, 32), &[255, 0, 0, 255], "cube in the center is red");
    assert_eq!(px(1, 1), &[255, 255, 255, 255], "corner is background");
}

//-------------------------------------------------
// camera selection (no GPU needed)

fn view(pos: [f32; 3], target: [f32; 3], up: [f32; 3], fov: f32) -> CameraView {
    CameraView { position: pos, target, up, fov_deg: fov }
}

#[test]
fn scene_camera_used_without_override() {
    let s = Scene::from_json(SCENE.as_bytes()).unwrap();
    let setup = camera_setup(&s, &Options::default()).unwrap();
    assert_eq!(setup.eye.to_array(), [0.0, 0.0, 5.0]);
    assert_eq!(setup.target.to_array(), [0.0, 0.0, 0.0]);
    assert_eq!(setup.fov_deg, 60.0);
}

#[test]
fn explicit_camera_overrides_scene_camera_and_fit() {
    let s = Scene::from_json(SCENE.as_bytes()).unwrap();
    let opts = Options {
        fit: true, // ignored when a camera is given
        camera: Some(view([3.0, 4.0, 5.0], [1.0, 0.0, 0.0], [0.0, 1.0, 0.0], 75.0)),
        ..Options::default()
    };
    let setup = camera_setup(&s, &opts).unwrap();
    assert_eq!(setup.eye.to_array(), [3.0, 4.0, 5.0]);
    assert_eq!(setup.target.to_array(), [1.0, 0.0, 0.0]);
    assert_eq!(setup.up.to_array(), [0.0, 1.0, 0.0]);
    assert_eq!(setup.fov_deg, 75.0);
}

#[test]
fn explicit_camera_up_is_kept_when_rolled() {
    // the trackball rotates `up`, so it is not always +Y. it must be passed through as-is
    let s = Scene::from_json(SCENE.as_bytes()).unwrap();
    let opts = Options { camera: Some(view([10.0, 10.0, 10.0], [0.0; 3], [0.3, 0.9, 0.1], 75.0)), ..Options::default() };
    let setup = camera_setup(&s, &opts).unwrap();
    assert_eq!(setup.up.to_array(), [0.3, 0.9, 0.1]);
}

#[test]
fn invalid_cameras_are_rejected() {
    let same_point = view([1.0, 1.0, 1.0], [1.0, 1.0, 1.0], [0.0, 1.0, 0.0], 75.0);
    let parallel_up = view([0.0, 0.0, 5.0], [0.0, 0.0, 0.0], [0.0, 0.0, 1.0], 75.0);
    let bad_fov = view([0.0, 0.0, 5.0], [0.0, 0.0, 0.0], [0.0, 1.0, 0.0], 180.0);
    let not_finite = view([f32::NAN, 0.0, 5.0], [0.0, 0.0, 0.0], [0.0, 1.0, 0.0], 75.0);
    let zero_up = view([0.0, 0.0, 5.0], [0.0, 0.0, 0.0], [0.0, 0.0, 0.0], 75.0);
    for bad in [same_point, parallel_up, bad_fov, not_finite, zero_up] {
        assert!(bad.validate().is_err(), "expected error for {bad:?}");
    }
    assert!(view([0.0, 0.0, 5.0], [0.0, 0.0, 0.0], [0.0, 1.0, 0.0], 75.0).validate().is_ok());
}
