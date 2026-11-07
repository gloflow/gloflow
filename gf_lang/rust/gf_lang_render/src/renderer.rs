//! Offscreen wgpu rendering of a [`Scene`].
//!
//! Mirrors the browser engine (ts/engine/gf_engine.ts, three.js):
//!  - cubes/spheres with flat-colored (unlit) `MeshBasicMaterial`, polygon-offset pushed back
//!  - optional black wireframe overlay (EdgesGeometry)
//!  - black 1px lines
//!  - `ShaderMaterial`s from user GLSL (see `shader.rs`)
//!  - 4x MSAA (three `antialias: true`)

use crate::geometry::{self, Mesh, Vertex};
use crate::scene::{Entity, Kind, Material, Scene};
use crate::shader::{self, Translated};
use bytemuck::{Pod, Zeroable};
use glam::{EulerRot, Mat3, Mat4, Vec3};
use std::collections::HashMap;
use wgpu::util::DeviceExt;

const COLOR_FORMAT: wgpu::TextureFormat = wgpu::TextureFormat::Rgba8Unorm;
const DEPTH_FORMAT: wgpu::TextureFormat = wgpu::TextureFormat::Depth32Float;

#[derive(Debug, Clone)]
pub struct Options {
    pub width: u32,
    pub height: u32,
    /// MSAA sample count: 1, 2 or 4
    pub samples: u32,
    /// ignore the scene camera position and frame all entities, looking from the scene camera's direction
    /// (the browser engine has interactive trackball controls, a static render has no such luxury)
    pub fit: bool,
    /// the view to render from, instead of the scene camera (ex. the view the user has in the dashboard).
    /// takes precedence over `fit`.
    pub camera: Option<CameraView>,
}

impl Default for Options {
    fn default() -> Self {
        Options { width: 1024, height: 768, samples: 4, fit: false, camera: None }
    }
}

/// An explicit camera: the same three.js-style values the browser engine has (position, look-at
/// target, up vector, vertical fov in degrees).
#[derive(Debug, Clone, Copy, PartialEq)]
pub struct CameraView {
    pub position: [f32; 3],
    pub target: [f32; 3],
    pub up: [f32; 3],
    pub fov_deg: f32,
}

impl CameraView {
    pub fn validate(&self) -> Result<(), String> {
        let all_finite = |v: &[f32; 3]| v.iter().all(|x| x.is_finite());
        if !(all_finite(&self.position) && all_finite(&self.target) && all_finite(&self.up)) {
            return Err("camera values must be finite numbers".into());
        }
        if !(self.fov_deg > 0.0 && self.fov_deg < 180.0) {
            return Err("camera fov_deg must be between 0 and 180".into());
        }
        let dir = Vec3::from(self.target) - Vec3::from(self.position);
        if dir.length_squared() < 1e-12 {
            return Err("camera position and target are the same point".into());
        }
        let up = Vec3::from(self.up).normalize_or_zero();
        if up == Vec3::ZERO {
            return Err("camera up must not be zero".into());
        }
        if dir.normalize().cross(up).length_squared() < 1e-12 {
            return Err("camera up must not be parallel to the view direction".into());
        }
        Ok(())
    }
}

/// Where the camera is and how it is oriented, for given options and scene.
#[derive(Debug, Clone, Copy, PartialEq)]
pub struct CameraSetup {
    pub eye: Vec3,
    pub target: Vec3,
    pub up: Vec3,
    pub fov_deg: f32,
}

/// The camera for `opts`: the explicit camera if given, else the scene camera (framed on the
/// entities with `fit`).
pub fn camera_setup(scene: &Scene, opts: &Options) -> Result<CameraSetup, String> {
    if let Some(c) = &opts.camera {
        c.validate()?;
        return Ok(CameraSetup {
            eye: Vec3::from(c.position),
            target: Vec3::from(c.target),
            up: Vec3::from(c.up),
            fov_deg: c.fov_deg,
        });
    }

    let cam = &scene.camera;
    let mut eye = Vec3::from(cam.position);
    let mut target = Vec3::from(cam.target);
    if opts.fit {
        if let Some((center, radius)) = scene_bounds(scene) {
            let dir = (eye - target).normalize_or(Vec3::ONE.normalize());
            let dist = radius / (cam.fov_deg.to_radians() * 0.5).sin().min((cam.fov_deg.to_radians() * 0.5 * opts.width as f32 / opts.height as f32).sin()) * 1.1;
            target = center;
            eye = center + dir * dist;
        }
    }
    let mut up = Vec3::Y;
    if (target - eye).normalize_or_zero().dot(up).abs() > 0.9999 {
        up = Vec3::Z;
    }
    Ok(CameraSetup { eye, target, up, fov_deg: cam.fov_deg })
}

pub struct Image {
    pub width: u32,
    pub height: u32,
    /// RGBA8, row-major, tightly packed
    pub rgba: Vec<u8>,
}

#[repr(C)]
#[derive(Clone, Copy, Pod, Zeroable)]
struct Instance {
    model: [[f32; 4]; 4],
    color: [f32; 4],
}

const BASIC_WGSL: &str = r#"
struct Cam { view_proj: mat4x4<f32> };
@group(0) @binding(0) var<uniform> cam: Cam;

struct VIn {
    @location(0) pos: vec3<f32>,
    @location(3) m0: vec4<f32>,
    @location(4) m1: vec4<f32>,
    @location(5) m2: vec4<f32>,
    @location(6) m3: vec4<f32>,
    @location(7) color: vec4<f32>,
};
struct VOut {
    @builtin(position) pos: vec4<f32>,
    @location(0) color: vec4<f32>,
};

@vertex
fn vs_main(in: VIn) -> VOut {
    let model = mat4x4<f32>(in.m0, in.m1, in.m2, in.m3);
    var out: VOut;
    out.pos = cam.view_proj * model * vec4<f32>(in.pos, 1.0);
    out.color = in.color;
    return out;
}

@fragment
fn fs_main(in: VOut) -> @location(0) vec4<f32> {
    return in.color;
}
"#;

fn err<E: std::fmt::Display>(ctx: &str) -> impl FnOnce(E) -> String + '_ {
    move |e| format!("{ctx}: {e}")
}

struct Buf {
    buf: wgpu::Buffer,
    count: u32,
}

fn vertex_buf(device: &wgpu::Device, label: &str, data: &[Vertex]) -> Option<Buf> {
    if data.is_empty() {
        return None;
    }
    Some(Buf {
        buf: device.create_buffer_init(&wgpu::util::BufferInitDescriptor {
            label: Some(label),
            contents: bytemuck::cast_slice(data),
            usage: wgpu::BufferUsages::VERTEX,
        }),
        count: data.len() as u32,
    })
}

fn instance_buf(device: &wgpu::Device, label: &str, data: &[Instance]) -> Option<Buf> {
    if data.is_empty() {
        return None;
    }
    Some(Buf {
        buf: device.create_buffer_init(&wgpu::util::BufferInitDescriptor {
            label: Some(label),
            contents: bytemuck::cast_slice(data),
            usage: wgpu::BufferUsages::VERTEX,
        }),
        count: data.len() as u32,
    })
}

fn model_matrix(e: &Entity) -> Mat4 {
    Mat4::from_translation(Vec3::from(e.position))
        * Mat4::from_euler(EulerRot::XYZ, e.rotation[0], e.rotation[1], e.rotation[2])
        * Mat4::from_scale(Vec3::from(e.scale))
}

fn instance(e: &Entity, color: [f32; 3]) -> Instance {
    Instance {
        model: model_matrix(e).to_cols_array_2d(),
        color: [color[0], color[1], color[2], 1.0],
    }
}

/// bounding sphere (center, radius) of all entities
fn scene_bounds(scene: &Scene) -> Option<(Vec3, f32)> {
    let mut min = Vec3::splat(f32::MAX);
    let mut max = Vec3::splat(f32::MIN);
    let mut any = false;
    let mut add = |p: Vec3, r: f32| {
        min = min.min(p - Vec3::splat(r));
        max = max.max(p + Vec3::splat(r));
        any = true;
    };
    for e in &scene.entities {
        let smax = e.scale.iter().fold(0.0f32, |a, b| a.max(b.abs()));
        match e.kind {
            Kind::Cube => add(Vec3::from(e.position), 0.87 * smax),
            Kind::Sphere => add(Vec3::from(e.position), geometry::SPHERE_RADIUS * smax),
            Kind::Line => e.points.iter().for_each(|p| add(Vec3::from(*p), 0.0)),
        }
    }
    any.then(|| ((min + max) * 0.5, ((max - min) * 0.5).length().max(1.0)))
}

/// right-handed view matrix (camera looks down -Z), like three.js `lookAt`
fn look_at_rh(eye: Vec3, target: Vec3, up: Vec3) -> Mat4 {
    let f = (target - eye).normalize();
    let s = f.cross(up).normalize();
    let u = s.cross(f);
    Mat4::from_cols(
        glam::Vec4::new(s.x, u.x, -f.x, 0.0),
        glam::Vec4::new(s.y, u.y, -f.y, 0.0),
        glam::Vec4::new(s.z, u.z, -f.z, 0.0),
        glam::Vec4::new(-s.dot(eye), -u.dot(eye), f.dot(eye), 1.0),
    )
}

/// right-handed perspective projection with wgpu clip space (depth 0..1)
fn perspective_rh(fov_y_rad: f32, aspect: f32, near: f32, far: f32) -> Mat4 {
    let f = 1.0 / (fov_y_rad * 0.5).tan();
    Mat4::from_cols(
        glam::Vec4::new(f / aspect, 0.0, 0.0, 0.0),
        glam::Vec4::new(0.0, f, 0.0, 0.0),
        glam::Vec4::new(0.0, 0.0, far / (near - far), -1.0),
        glam::Vec4::new(0.0, 0.0, near * far / (near - far), 0.0),
    )
}

const BLACK: [f32; 3] = [0.0, 0.0, 0.0];

/// Renders the scene offscreen and returns an RGBA8 image.
pub fn render(scene: &Scene, opts: &Options) -> Result<Image, String> {
    if opts.width == 0 || opts.height == 0 {
        return Err("width and height must be > 0".into());
    }
    if !matches!(opts.samples, 1 | 2 | 4) {
        return Err("samples must be 1, 2 or 4".into());
    }

    //------------------------
    // DEVICE
    let instance_gpu = wgpu::Instance::new(wgpu::InstanceDescriptor::new_without_display_handle());

    let adapter = pollster::block_on(instance_gpu.request_adapter(&wgpu::RequestAdapterOptions {
        power_preference: wgpu::PowerPreference::HighPerformance,
        force_fallback_adapter: false,
        compatible_surface: None,
        apply_limit_buckets: false,
    }))
    .or_else(|_| {
        pollster::block_on(instance_gpu.request_adapter(&wgpu::RequestAdapterOptions {
            power_preference: wgpu::PowerPreference::LowPower,
            force_fallback_adapter: true,
            compatible_surface: None,
            apply_limit_buckets: false,
        }))
    })
    .map_err(|e| {
        format!("no GPU adapter found ({e}). On servers without a GPU install a software Vulkan driver (Mesa lavapipe: mesa-vulkan-drivers)")
    })?;

    let (device, queue) = pollster::block_on(adapter.request_device(&wgpu::DeviceDescriptor {
        label: Some("gf_lang_render"),
        ..Default::default()
    }))
    .map_err(err("request_device"))?;

    //------------------------
    // CAMERA
    let setup = camera_setup(scene, opts)?;
    let view = look_at_rh(setup.eye, setup.target, setup.up);
    let aspect = opts.width as f32 / opts.height as f32;
    let proj = perspective_rh(setup.fov_deg.to_radians(), aspect, scene.camera.near, scene.camera.far);
    let view_proj = proj * view;

    let cam_buf = device.create_buffer_init(&wgpu::util::BufferInitDescriptor {
        label: Some("camera"),
        contents: bytemuck::cast_slice(&view_proj.to_cols_array()),
        usage: wgpu::BufferUsages::UNIFORM,
    });

    //------------------------
    // BASIC PIPELINES
    let basic_module = device.create_shader_module(wgpu::ShaderModuleDescriptor {
        label: Some("gf_basic"),
        source: wgpu::ShaderSource::Wgsl(BASIC_WGSL.into()),
    });
    let cam_bgl = device.create_bind_group_layout(&wgpu::BindGroupLayoutDescriptor {
        label: Some("camera_bgl"),
        entries: &[wgpu::BindGroupLayoutEntry {
            binding: 0,
            visibility: wgpu::ShaderStages::VERTEX,
            ty: wgpu::BindingType::Buffer {
                ty: wgpu::BufferBindingType::Uniform,
                has_dynamic_offset: false,
                min_binding_size: None,
            },
            count: None,
        }],
    });
    let cam_bg = device.create_bind_group(&wgpu::BindGroupDescriptor {
        label: Some("camera_bg"),
        layout: &cam_bgl,
        entries: &[wgpu::BindGroupEntry { binding: 0, resource: cam_buf.as_entire_binding() }],
    });
    let basic_layout = device.create_pipeline_layout(&wgpu::PipelineLayoutDescriptor {
        label: Some("basic_layout"),
        bind_group_layouts: &[Some(&cam_bgl)],
        immediate_size: 0,
    });

    let vertex_attrs = wgpu::vertex_attr_array![0 => Float32x3, 1 => Float32x3, 2 => Float32x2];
    let vertex_layout = wgpu::VertexBufferLayout {
        array_stride: std::mem::size_of::<Vertex>() as u64,
        step_mode: wgpu::VertexStepMode::Vertex,
        attributes: &vertex_attrs,
    };
    let inst_attrs = wgpu::vertex_attr_array![3 => Float32x4, 4 => Float32x4, 5 => Float32x4, 6 => Float32x4, 7 => Float32x4];
    let inst_layout = wgpu::VertexBufferLayout {
        array_stride: std::mem::size_of::<Instance>() as u64,
        step_mode: wgpu::VertexStepMode::Instance,
        attributes: &inst_attrs,
    };

    let multisample = wgpu::MultisampleState { count: opts.samples, mask: !0, alpha_to_coverage_enabled: false };

    // filled geometry is pushed back (three.js polygonOffset) so wireframes drawn over it are visible
    let make_depth = |bias: wgpu::DepthBiasState| wgpu::DepthStencilState {
        format: DEPTH_FORMAT,
        depth_write_enabled: Some(true),
        depth_compare: Some(wgpu::CompareFunction::Less),
        stencil: wgpu::StencilState::default(),
        bias,
    };
    let pushed_back = wgpu::DepthBiasState { constant: 1, slope_scale: 1.0, clamp: 0.0 };

    let make_basic_pipeline = |label: &str, topology: wgpu::PrimitiveTopology, bias: wgpu::DepthBiasState| {
        device.create_render_pipeline(&wgpu::RenderPipelineDescriptor {
            label: Some(label),
            layout: Some(&basic_layout),
            vertex: wgpu::VertexState {
                module: &basic_module,
                entry_point: Some("vs_main"),
                compilation_options: Default::default(),
                buffers: &[Some(vertex_layout.clone()), Some(inst_layout.clone())],
            },
            primitive: wgpu::PrimitiveState { topology, cull_mode: None, ..Default::default() },
            depth_stencil: Some(make_depth(bias)),
            multisample,
            fragment: Some(wgpu::FragmentState {
                module: &basic_module,
                entry_point: Some("fs_main"),
                compilation_options: Default::default(),
                targets: &[Some(COLOR_FORMAT.into())],
            }),
            multiview_mask: None,
            cache: None,
        })
    };
    let tri_pipeline = make_basic_pipeline("basic_tris", wgpu::PrimitiveTopology::TriangleList, pushed_back);
    let line_pipeline = make_basic_pipeline("basic_lines", wgpu::PrimitiveTopology::LineList, wgpu::DepthBiasState::default());

    //------------------------
    // MESHES + INSTANCES
    let cube_mesh = geometry::cube();
    let sphere_mesh = geometry::sphere();
    let meshes: [(Kind, &Mesh); 2] = [(Kind::Cube, &cube_mesh), (Kind::Sphere, &sphere_mesh)];

    struct Draw {
        tris: Option<Buf>,
        edges: Option<Buf>,
        basic_inst: Option<Buf>,
        edge_inst: Option<Buf>,
    }
    let mut draws: Vec<Draw> = Vec::new();
    for (kind, mesh) in meshes {
        let mut basic = Vec::new();
        let mut edge = Vec::new();
        for e in scene.entities.iter().filter(|e| e.kind == kind) {
            if matches!(e.material, Material::Basic) {
                basic.push(instance(e, e.color));
            }
            if e.wireframe {
                edge.push(instance(e, BLACK));
            }
        }
        draws.push(Draw {
            tris: vertex_buf(&device, "mesh_tris", &mesh.tris),
            edges: vertex_buf(&device, "mesh_edges", &mesh.edges),
            basic_inst: instance_buf(&device, "basic_inst", &basic),
            edge_inst: instance_buf(&device, "edge_inst", &edge),
        });
    }

    // lines
    let mut line_verts: Vec<Vertex> = Vec::new();
    for e in scene.entities.iter().filter(|e| e.kind == Kind::Line) {
        if e.points.len() >= 2 {
            for w in e.points.windows(2) {
                for p in w {
                    line_verts.push(Vertex { position: *p, normal: [0.0; 3], uv: [0.0; 2] });
                }
            }
        }
    }
    let line_buf = vertex_buf(&device, "lines", &line_verts);
    let identity_inst = instance_buf(
        &device,
        "identity_inst",
        &[Instance { model: Mat4::IDENTITY.to_cols_array_2d(), color: [0.0, 0.0, 0.0, 1.0] }],
    );

    //------------------------
    // SHADER (custom GLSL) PIPELINES
    struct ShaderPipeline {
        translated: Translated,
        pipeline: wgpu::RenderPipeline,
        bgl: wgpu::BindGroupLayout,
    }
    let mut shader_pipelines: HashMap<&str, ShaderPipeline> = HashMap::new();
    for e in &scene.entities {
        let Material::Shader { shader: name, .. } = &e.material else { continue };
        if shader_pipelines.contains_key(name.as_str()) {
            continue;
        }
        let def = scene.shaders.get(name).ok_or_else(|| format!("unknown shader '{name}'"))?;
        let translated = shader::translate(def).map_err(|m| format!("shader '{name}': {m}"))?;

        let make_module = |stage: wgpu::naga::ShaderStage, src: &str| -> Result<wgpu::ShaderModule, String> {
            let scope = device.push_error_scope(wgpu::ErrorFilter::Validation);
            let module = device.create_shader_module(wgpu::ShaderModuleDescriptor {
                label: Some(name),
                source: wgpu::ShaderSource::Glsl {
                    shader: src.to_string().into(),
                    stage,
                    defines: &[],
                },
            });
            match pollster::block_on(scope.pop()) {
                Some(e) => Err(format!("shader '{name}' {stage:?} stage failed to compile: {e}\n--- translated source ---\n{src}")),
                None => Ok(module),
            }
        };
        let vs = make_module(wgpu::naga::ShaderStage::Vertex, &translated.vertex)?;
        let fs = make_module(wgpu::naga::ShaderStage::Fragment, &translated.fragment)?;

        let mut entries = vec![wgpu::BindGroupLayoutEntry {
            binding: 0,
            visibility: wgpu::ShaderStages::VERTEX_FRAGMENT,
            ty: wgpu::BindingType::Buffer {
                ty: wgpu::BufferBindingType::Uniform,
                has_dynamic_offset: false,
                min_binding_size: None,
            },
            count: None,
        }];
        if translated.block_size > 0 {
            entries.push(wgpu::BindGroupLayoutEntry {
                binding: 1,
                visibility: wgpu::ShaderStages::VERTEX_FRAGMENT,
                ty: wgpu::BindingType::Buffer {
                    ty: wgpu::BufferBindingType::Uniform,
                    has_dynamic_offset: false,
                    min_binding_size: None,
                },
                count: None,
            });
        }
        let bgl = device.create_bind_group_layout(&wgpu::BindGroupLayoutDescriptor {
            label: Some("shader_bgl"),
            entries: &entries,
        });
        let layout = device.create_pipeline_layout(&wgpu::PipelineLayoutDescriptor {
            label: Some("shader_layout"),
            bind_group_layouts: &[Some(&bgl)],
            immediate_size: 0,
        });

        let scope = device.push_error_scope(wgpu::ErrorFilter::Validation);
        let pipeline = device.create_render_pipeline(&wgpu::RenderPipelineDescriptor {
            label: Some(name),
            layout: Some(&layout),
            vertex: wgpu::VertexState {
                module: &vs,
                entry_point: Some("main"),
                compilation_options: Default::default(),
                buffers: &[Some(vertex_layout.clone())],
            },
            primitive: wgpu::PrimitiveState { cull_mode: None, ..Default::default() },
            depth_stencil: Some(make_depth(wgpu::DepthBiasState::default())),
            multisample,
            fragment: Some(wgpu::FragmentState {
                module: &fs,
                entry_point: Some("main"),
                compilation_options: Default::default(),
                targets: &[Some(COLOR_FORMAT.into())],
            }),
            multiview_mask: None,
            cache: None,
        });
        if let Some(e) = pollster::block_on(scope.pop()) {
            return Err(format!(
                "shader '{name}' pipeline creation failed: {e}\n--- vertex ---\n{}\n--- fragment ---\n{}",
                translated.vertex, translated.fragment
            ));
        }
        shader_pipelines.insert(name.as_str(), ShaderPipeline { translated, pipeline, bgl });
    }

    // per-entity shader draws: (pipeline name, mesh index, bind group)
    let mut shader_draws: Vec<(&str, usize, wgpu::BindGroup)> = Vec::new();
    for e in &scene.entities {
        let Material::Shader { shader: name, uniforms } = &e.material else { continue };
        let mesh_idx = match e.kind {
            Kind::Cube => 0,
            Kind::Sphere => 1,
            Kind::Line => continue,
        };
        let sp = &shader_pipelines[name.as_str()];

        let model = model_matrix(e);
        let model_view = view * model;
        let normal_m = Mat3::from_mat4(model_view).inverse().transpose();

        let mut b: Vec<f32> = Vec::with_capacity(shader::BUILTIN_BLOCK_SIZE / 4);
        b.extend_from_slice(&proj.to_cols_array());
        b.extend_from_slice(&model_view.to_cols_array());
        b.extend_from_slice(&view.to_cols_array());
        b.extend_from_slice(&model.to_cols_array());
        for c in [normal_m.x_axis, normal_m.y_axis, normal_m.z_axis] {
            b.extend_from_slice(&[c.x, c.y, c.z, 0.0]);
        }
        b.extend_from_slice(&[setup.eye.x, setup.eye.y, setup.eye.z, 0.0]);
        debug_assert_eq!(b.len() * 4, shader::BUILTIN_BLOCK_SIZE);

        let builtin_buf = device.create_buffer_init(&wgpu::util::BufferInitDescriptor {
            label: Some("builtin_uniforms"),
            contents: bytemuck::cast_slice(&b),
            usage: wgpu::BufferUsages::UNIFORM,
        });
        let user_buf = if sp.translated.block_size > 0 {
            let bytes = shader::pack_uniforms(&sp.translated, uniforms).map_err(|m| format!("shader '{name}': {m}"))?;
            Some(device.create_buffer_init(&wgpu::util::BufferInitDescriptor {
                label: Some("user_uniforms"),
                contents: &bytes,
                usage: wgpu::BufferUsages::UNIFORM,
            }))
        } else {
            None
        };
        let mut entries = vec![wgpu::BindGroupEntry { binding: 0, resource: builtin_buf.as_entire_binding() }];
        if let Some(ub) = &user_buf {
            entries.push(wgpu::BindGroupEntry { binding: 1, resource: ub.as_entire_binding() });
        }
        let bg = device.create_bind_group(&wgpu::BindGroupDescriptor {
            label: Some("shader_bg"),
            layout: &sp.bgl,
            entries: &entries,
        });
        shader_draws.push((name.as_str(), mesh_idx, bg));
    }

    //------------------------
    // TARGETS
    let size = wgpu::Extent3d { width: opts.width, height: opts.height, depth_or_array_layers: 1 };
    let resolve_tex = device.create_texture(&wgpu::TextureDescriptor {
        label: Some("resolve"),
        size,
        mip_level_count: 1,
        sample_count: 1,
        dimension: wgpu::TextureDimension::D2,
        format: COLOR_FORMAT,
        usage: wgpu::TextureUsages::RENDER_ATTACHMENT | wgpu::TextureUsages::COPY_SRC,
        view_formats: &[],
    });
    let resolve_view = resolve_tex.create_view(&Default::default());
    let msaa_view = (opts.samples > 1).then(|| {
        device
            .create_texture(&wgpu::TextureDescriptor {
                label: Some("msaa"),
                size,
                mip_level_count: 1,
                sample_count: opts.samples,
                dimension: wgpu::TextureDimension::D2,
                format: COLOR_FORMAT,
                usage: wgpu::TextureUsages::RENDER_ATTACHMENT,
                view_formats: &[],
            })
            .create_view(&Default::default())
    });
    let depth_view = device
        .create_texture(&wgpu::TextureDescriptor {
            label: Some("depth"),
            size,
            mip_level_count: 1,
            sample_count: opts.samples,
            dimension: wgpu::TextureDimension::D2,
            format: DEPTH_FORMAT,
            usage: wgpu::TextureUsages::RENDER_ATTACHMENT,
            view_formats: &[],
        })
        .create_view(&Default::default());

    let bpr = (opts.width * 4).div_ceil(wgpu::COPY_BYTES_PER_ROW_ALIGNMENT) * wgpu::COPY_BYTES_PER_ROW_ALIGNMENT;
    let readback = device.create_buffer(&wgpu::BufferDescriptor {
        label: Some("readback"),
        size: (bpr * opts.height) as u64,
        usage: wgpu::BufferUsages::COPY_DST | wgpu::BufferUsages::MAP_READ,
        mapped_at_creation: false,
    });

    //------------------------
    // DRAW
    let bg = scene.background.unwrap_or([0.0; 3]);
    let clear = wgpu::Color { r: bg[0] as f64, g: bg[1] as f64, b: bg[2] as f64, a: 1.0 };

    let mut encoder = device.create_command_encoder(&Default::default());
    {
        let (view_att, resolve_att) = match &msaa_view {
            Some(m) => (m, Some(&resolve_view)),
            None => (&resolve_view, None),
        };
        let mut pass = encoder.begin_render_pass(&wgpu::RenderPassDescriptor {
            label: Some("scene"),
            color_attachments: &[Some(wgpu::RenderPassColorAttachment {
                view: view_att,
                depth_slice: None,
                resolve_target: resolve_att,
                ops: wgpu::Operations { load: wgpu::LoadOp::Clear(clear), store: wgpu::StoreOp::Store },
            })],
            depth_stencil_attachment: Some(wgpu::RenderPassDepthStencilAttachment {
                view: &depth_view,
                depth_ops: Some(wgpu::Operations { load: wgpu::LoadOp::Clear(1.0), store: wgpu::StoreOp::Discard }),
                stencil_ops: None,
            }),
            timestamp_writes: None,
            occlusion_query_set: None,
            multiview_mask: None,
        });

        // basic filled meshes
        pass.set_pipeline(&tri_pipeline);
        pass.set_bind_group(0, &cam_bg, &[]);
        for d in &draws {
            if let (Some(t), Some(i)) = (&d.tris, &d.basic_inst) {
                pass.set_vertex_buffer(0, t.buf.slice(..));
                pass.set_vertex_buffer(1, i.buf.slice(..));
                pass.draw(0..t.count, 0..i.count);
            }
        }

        // custom shader meshes
        for (name, mesh_idx, bind_group) in &shader_draws {
            if let Some(t) = &draws[*mesh_idx].tris {
                pass.set_pipeline(&shader_pipelines[name].pipeline);
                pass.set_bind_group(0, bind_group, &[]);
                pass.set_vertex_buffer(0, t.buf.slice(..));
                pass.draw(0..t.count, 0..1);
            }
        }

        // wireframes + lines
        pass.set_pipeline(&line_pipeline);
        pass.set_bind_group(0, &cam_bg, &[]);
        for d in &draws {
            if let (Some(e), Some(i)) = (&d.edges, &d.edge_inst) {
                pass.set_vertex_buffer(0, e.buf.slice(..));
                pass.set_vertex_buffer(1, i.buf.slice(..));
                pass.draw(0..e.count, 0..i.count);
            }
        }
        if let (Some(l), Some(i)) = (&line_buf, &identity_inst) {
            pass.set_vertex_buffer(0, l.buf.slice(..));
            pass.set_vertex_buffer(1, i.buf.slice(..));
            pass.draw(0..l.count, 0..1);
        }
    }

    encoder.copy_texture_to_buffer(
        wgpu::TexelCopyTextureInfo {
            texture: &resolve_tex,
            mip_level: 0,
            origin: wgpu::Origin3d::ZERO,
            aspect: wgpu::TextureAspect::All,
        },
        wgpu::TexelCopyBufferInfo {
            buffer: &readback,
            layout: wgpu::TexelCopyBufferLayout { offset: 0, bytes_per_row: Some(bpr), rows_per_image: Some(opts.height) },
        },
        size,
    );
    queue.submit([encoder.finish()]);

    //------------------------
    // READBACK
    let slice = readback.slice(..);
    let (tx, rx) = std::sync::mpsc::channel();
    slice.map_async(wgpu::MapMode::Read, move |r| {
        let _ = tx.send(r);
    });
    device.poll(wgpu::PollType::wait_indefinitely()).map_err(err("device poll"))?;
    rx.recv().map_err(err("readback channel"))?.map_err(err("buffer map"))?;

    let data = slice.get_mapped_range().map_err(err("map range"))?;
    let row = (opts.width * 4) as usize;
    let mut rgba = Vec::with_capacity(row * opts.height as usize);
    for y in 0..opts.height as usize {
        rgba.extend_from_slice(&data[y * bpr as usize..y * bpr as usize + row]);
    }
    drop(data);
    readback.unmap();

    Ok(Image { width: opts.width, height: opts.height, rgba })
}
