//! Translation of three.js-style `ShaderMaterial` GLSL (GLSL ES 1.00, as embedded in
//! gf_lang programs) to Vulkan-style GLSL 450, which naga/wgpu can consume.
//!
//! three.js prepends built-in declarations to user shaders, so user code relies on
//! them without declaring them. We do the same:
//!
//!   uniforms:   projectionMatrix modelViewMatrix viewMatrix modelMatrix normalMatrix cameraPosition
//!   attributes: position normal uv (vertex shader only)
//!
//! and rewrite `varying`/`attribute`/`uniform` declarations, `gl_FragColor`, and drop
//! `precision`/`#version` lines.
//!
//! Bind group 0: binding 0 = built-ins block, binding 1 = user uniforms block (only if any).

use crate::scene::ShaderDef;
use serde_json::Value;

pub const BUILTIN_BLOCK_SIZE: usize = 320; // 4 mat4 + mat3 (3 padded cols) + vec3 (padded)

#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub enum UType {
    Float,
    Int,
    Vec2,
    Vec3,
    Vec4,
}

impl UType {
    fn parse(s: &str) -> Option<UType> {
        Some(match s {
            "float" => UType::Float,
            "int" => UType::Int,
            "vec2" => UType::Vec2,
            "vec3" => UType::Vec3,
            "vec4" => UType::Vec4,
            _ => return None,
        })
    }
    fn glsl(self) -> &'static str {
        match self {
            UType::Float => "float",
            UType::Int => "int",
            UType::Vec2 => "vec2",
            UType::Vec3 => "vec3",
            UType::Vec4 => "vec4",
        }
    }
    /// (align, size) in std140
    fn layout(self) -> (usize, usize) {
        match self {
            UType::Float | UType::Int => (4, 4),
            UType::Vec2 => (8, 8),
            UType::Vec3 => (16, 12),
            UType::Vec4 => (16, 16),
        }
    }
}

#[derive(Debug, Clone)]
pub struct UniformSlot {
    pub name: String,
    pub ty: UType,
    pub offset: usize,
    pub default: Value,
}

pub struct Translated {
    pub vertex: String,
    pub fragment: String,
    pub uniforms: Vec<UniformSlot>,
    /// size of the user uniforms block, 0 if there are none
    pub block_size: usize,
}

#[derive(Default)]
struct Parsed {
    body: String,
    varyings: Vec<(String, String)>, // (type, name) in declaration order
    uniforms: Vec<(String, String)>, // (type, name)
}

const BUILTIN_UNIFORMS: [&str; 6] = [
    "projectionMatrix",
    "modelViewMatrix",
    "viewMatrix",
    "modelMatrix",
    "normalMatrix",
    "cameraPosition",
];

fn strip_comments(src: &str) -> String {
    let b = src.as_bytes();
    let mut out = String::with_capacity(src.len());
    let mut i = 0;
    while i < b.len() {
        if b[i] == b'/' && i + 1 < b.len() && b[i + 1] == b'/' {
            while i < b.len() && b[i] != b'\n' {
                i += 1;
            }
        } else if b[i] == b'/' && i + 1 < b.len() && b[i + 1] == b'*' {
            i += 2;
            while i + 1 < b.len() && !(b[i] == b'*' && b[i + 1] == b'/') {
                if b[i] == b'\n' {
                    out.push('\n');
                }
                i += 1;
            }
            i += 2;
        } else {
            // copy a whole utf8 char
            let ch = src[i..].chars().next().unwrap();
            out.push(ch);
            i += ch.len_utf8();
        }
    }
    out
}

fn parse_source(src: &str) -> Result<Parsed, String> {
    let src = strip_comments(src);
    let mut p = Parsed::default();

    for line in src.lines() {
        let t = line.trim_start();
        if t.starts_with("#version") || t.starts_with("precision ") {
            continue;
        }
        let kw = ["varying ", "uniform ", "attribute "]
            .iter()
            .find(|k| t.starts_with(**k))
            .copied();
        let Some(kw) = kw else {
            p.body.push_str(line);
            p.body.push('\n');
            continue;
        };

        // declaration statement up to the first ';', keep the rest of the line
        let Some(end) = t.find(';') else {
            return Err(format!("unsupported multi-line declaration: {t}"));
        };
        let decl = &t[kw.len()..end];
        let rest = &t[end + 1..];

        let mut toks = decl.split_whitespace().peekable();
        if matches!(toks.peek(), Some(&"lowp") | Some(&"mediump") | Some(&"highp")) {
            toks.next();
        }
        let ty = toks.next().ok_or_else(|| format!("bad declaration: {t}"))?.to_string();
        let names_str: String = toks.collect::<Vec<_>>().join(" ");
        for name in names_str.split(',') {
            let name = name.trim().to_string();
            if name.is_empty() {
                continue;
            }
            if name.contains('[') {
                return Err(format!("array declarations are not supported: {t}"));
            }
            match kw {
                "varying " => p.varyings.push((ty.clone(), name)),
                "uniform " => {
                    if !BUILTIN_UNIFORMS.contains(&name.as_str()) {
                        p.uniforms.push((ty.clone(), name));
                    }
                }
                _ => {} // attribute: position/normal/uv are declared by us
            }
        }
        if !rest.trim().is_empty() {
            p.body.push_str(rest);
            p.body.push('\n');
        }
    }
    Ok(p)
}

fn builtins_block() -> &'static str {
    "layout(set = 0, binding = 0) uniform GFbuiltin {\n\
     mat4 projectionMatrix;\n\
     mat4 modelViewMatrix;\n\
     mat4 viewMatrix;\n\
     mat4 modelMatrix;\n\
     mat3 normalMatrix;\n\
     vec3 cameraPosition;\n\
     };\n"
}

pub fn translate(def: &ShaderDef) -> Result<Translated, String> {
    let vs = parse_source(&def.vertex)?;
    let fs = parse_source(&def.fragment)?;

    //------------------------
    // USER UNIFORMS: union of declared (scene) and in-source uniforms
    let mut decls: Vec<(String, UType, Value)> = Vec::new();
    for u in &def.uniforms {
        let ty = UType::parse(&u.ty).ok_or_else(|| format!("unsupported uniform type '{}' ({})", u.ty, u.name))?;
        decls.push((u.name.clone(), ty, u.default.clone()));
    }
    for (ty, name) in vs.uniforms.iter().chain(fs.uniforms.iter()) {
        if decls.iter().any(|d| &d.0 == name) {
            continue;
        }
        let ty = UType::parse(ty).ok_or_else(|| format!("unsupported uniform type '{ty}' ({name})"))?;
        decls.push((name.clone(), ty, Value::Null));
    }

    let mut slots = Vec::new();
    let mut cur = 0usize;
    for (name, ty, default) in decls {
        let (align, size) = ty.layout();
        let offset = cur.div_ceil(align) * align;
        cur = offset + size;
        slots.push(UniformSlot { name, ty, offset, default });
    }
    let block_size = if slots.is_empty() { 0 } else { cur.div_ceil(16) * 16 };

    let mut user_block = String::new();
    if !slots.is_empty() {
        user_block.push_str("layout(set = 0, binding = 1) uniform GFuser {\n");
        for s in &slots {
            user_block.push_str(&format!("{} {};\n", s.ty.glsl(), s.name));
        }
        user_block.push_str("};\n");
    }

    //------------------------
    // VERTEX
    let mut vertex = String::from("#version 450\n");
    vertex.push_str(builtins_block());
    vertex.push_str(&user_block);
    vertex.push_str(
        "layout(location = 0) in vec3 position;\n\
         layout(location = 1) in vec3 normal;\n\
         layout(location = 2) in vec2 uv;\n",
    );
    for (i, (ty, name)) in vs.varyings.iter().enumerate() {
        vertex.push_str(&format!("layout(location = {i}) out {ty} {name};\n"));
    }
    vertex.push_str(&vs.body);

    //------------------------
    // FRAGMENT
    let mut fragment = String::from("#version 450\n");
    fragment.push_str(builtins_block());
    fragment.push_str(&user_block);
    for (ty, name) in &fs.varyings {
        let loc = vs
            .varyings
            .iter()
            .position(|(_, n)| n == name)
            .ok_or_else(|| format!("fragment varying '{name}' is not written by the vertex shader"))?;
        fragment.push_str(&format!("layout(location = {loc}) in {ty} {name};\n"));
    }
    fragment.push_str("layout(location = 0) out vec4 gf_FragColor;\n");
    fragment.push_str(&fs.body.replace("gl_FragColor", "gf_FragColor"));

    Ok(Translated { vertex, fragment, uniforms: slots, block_size })
}

/// Packs uniform values (scene entity overrides, then shader defaults, then zero) as std140.
pub fn pack_uniforms(
    t: &Translated,
    values: &std::collections::HashMap<String, Value>,
) -> Result<Vec<u8>, String> {
    let mut buf = vec![0u8; t.block_size];
    for s in &t.uniforms {
        let val = values.get(&s.name).unwrap_or(&s.default);
        let comps = components(val, s.ty).map_err(|e| format!("uniform '{}': {e}", s.name))?;
        for (i, c) in comps.iter().enumerate() {
            let o = s.offset + i * 4;
            match s.ty {
                UType::Int => buf[o..o + 4].copy_from_slice(&(*c as i32).to_le_bytes()),
                _ => buf[o..o + 4].copy_from_slice(&(*c as f32).to_le_bytes()),
            }
        }
    }
    Ok(buf)
}

fn components(val: &Value, ty: UType) -> Result<Vec<f64>, String> {
    let n = match ty {
        UType::Float | UType::Int => 1,
        UType::Vec2 => 2,
        UType::Vec3 => 3,
        UType::Vec4 => 4,
    };
    match val {
        Value::Null => Ok(vec![0.0; n]),
        Value::Number(x) if n == 1 => Ok(vec![x.as_f64().unwrap_or(0.0)]),
        Value::Bool(b) if n == 1 => Ok(vec![if *b { 1.0 } else { 0.0 }]),
        Value::Array(a) if a.len() == n => a
            .iter()
            .map(|x| x.as_f64().ok_or_else(|| "non-numeric component".to_string()))
            .collect(),
        other => Err(format!("value {other} does not match the {n}-component uniform type")),
    }
}
