//! Scene description, produced by the Go interpreter (`gf_lang/go/gf_lang/gf_scene.go`).

use serde::Deserialize;
use std::collections::HashMap;

pub const SUPPORTED_VERSION: u32 = 1;

#[derive(Debug, Deserialize)]
pub struct Scene {
    pub version: u32,
    /// [r,g,b] in 0..1. None == black (browser engine default).
    #[serde(default)]
    pub background: Option<[f32; 3]>,
    pub camera: Camera,
    #[serde(default)]
    pub shaders: HashMap<String, ShaderDef>,
    #[serde(default)]
    pub entities: Vec<Entity>,
    #[serde(default)]
    pub truncated: bool,
}

#[derive(Debug, Deserialize)]
pub struct Camera {
    pub position: [f32; 3],
    pub target: [f32; 3],
    pub fov_deg: f32,
    pub near: f32,
    pub far: f32,
}

#[derive(Debug, Deserialize)]
pub struct ShaderDef {
    #[serde(default)]
    pub uniforms: Vec<UniformDef>,
    pub vertex: String,
    pub fragment: String,
}

#[derive(Debug, Deserialize)]
pub struct UniformDef {
    pub name: String,
    #[serde(rename = "type")]
    pub ty: String,
    #[serde(default)]
    pub default: serde_json::Value,
}

#[derive(Debug, Clone, Copy, PartialEq, Eq, Deserialize)]
#[serde(rename_all = "lowercase")]
pub enum Kind {
    Cube,
    Sphere,
    Line,
}

#[derive(Debug, Deserialize)]
pub struct Entity {
    pub kind: Kind,
    #[serde(default)]
    pub position: [f32; 3],
    /// euler XYZ, radians
    #[serde(default)]
    pub rotation: [f32; 3],
    #[serde(default = "one3")]
    pub scale: [f32; 3],
    #[serde(default)]
    pub color: [f32; 3],
    #[serde(default)]
    pub wireframe: bool,
    pub material: Material,
    /// lines only, world-space
    #[serde(default)]
    pub points: Vec<[f32; 3]>,
    // "animation" is ignored for now (stills only).
}

fn one3() -> [f32; 3] {
    [1.0; 3]
}

#[derive(Debug, Deserialize)]
#[serde(tag = "type", rename_all = "lowercase")]
pub enum Material {
    Basic,
    Shader {
        shader: String,
        #[serde(default)]
        uniforms: HashMap<String, serde_json::Value>,
    },
}

impl Scene {
    pub fn from_json(bytes: &[u8]) -> Result<Scene, String> {
        let scene: Scene = serde_json::from_slice(bytes).map_err(|e| format!("invalid scene JSON: {e}"))?;
        if scene.version != SUPPORTED_VERSION {
            return Err(format!(
                "unsupported scene version {} (renderer supports {})",
                scene.version, SUPPORTED_VERSION
            ));
        }
        Ok(scene)
    }
}
