//! Headless wgpu renderer for gf_lang scenes.
//!
//! The gf_lang interpreter (Go) runs a program with a scene-recording extern API and
//! produces a [`scene::Scene`] (JSON). This crate renders that scene to an RGBA image
//! without a window or display, so it can run on a server (Vulkan, or software
//! Vulkan via Mesa lavapipe on machines without a GPU).
//!
//! Scene format docs: `gf_lang/docs/scene.md`.

pub mod geometry;
pub mod renderer;
pub mod scene;
pub mod shader;

pub use renderer::{camera_setup, render, CameraSetup, CameraView, Image, Options};
pub use scene::Scene;

/// Encodes an RGBA8 image as PNG bytes.
pub fn encode_png(image: &Image) -> Result<Vec<u8>, String> {
    let mut out = Vec::new();
    {
        let mut enc = png::Encoder::new(&mut out, image.width, image.height);
        enc.set_color(png::ColorType::Rgba);
        enc.set_depth(png::BitDepth::Eight);
        let mut w = enc.write_header().map_err(|e| e.to_string())?;
        w.write_image_data(&image.rgba).map_err(|e| e.to_string())?;
    }
    Ok(out)
}
