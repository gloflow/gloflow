//! Mesh generation, matching the three.js geometries used by the browser engine
//! (ts/engine/gf_engine.ts): BoxGeometry(1,1,1) and SphereGeometry(10, 50, 20),
//! plus their edges (used for the wireframe overlay).

use bytemuck::{Pod, Zeroable};

#[repr(C)]
#[derive(Clone, Copy, Pod, Zeroable)]
pub struct Vertex {
    pub position: [f32; 3],
    pub normal: [f32; 3],
    pub uv: [f32; 2],
}

pub const SPHERE_RADIUS: f32 = 10.0;
const SPHERE_WIDTH_SEGMENTS: usize = 50;
const SPHERE_HEIGHT_SEGMENTS: usize = 20;

pub struct Mesh {
    /// triangle list
    pub tris: Vec<Vertex>,
    /// line list (pairs of vertices)
    pub edges: Vec<Vertex>,
}

fn v(p: [f32; 3], n: [f32; 3], uv: [f32; 2]) -> Vertex {
    Vertex { position: p, normal: n, uv }
}

/// Unit cube centered at the origin.
pub fn cube() -> Mesh {
    // (normal, u axis, v axis): face corners = n*0.5 +- u*0.5 +- v*0.5, u x v == n (CCW seen from outside)
    let faces: [([f32; 3], [f32; 3], [f32; 3]); 6] = [
        ([1.0, 0.0, 0.0], [0.0, 1.0, 0.0], [0.0, 0.0, 1.0]),
        ([-1.0, 0.0, 0.0], [0.0, 0.0, 1.0], [0.0, 1.0, 0.0]),
        ([0.0, 1.0, 0.0], [0.0, 0.0, 1.0], [1.0, 0.0, 0.0]),
        ([0.0, -1.0, 0.0], [1.0, 0.0, 0.0], [0.0, 0.0, 1.0]),
        ([0.0, 0.0, 1.0], [1.0, 0.0, 0.0], [0.0, 1.0, 0.0]),
        ([0.0, 0.0, -1.0], [0.0, 1.0, 0.0], [1.0, 0.0, 0.0]),
    ];
    let mut tris = Vec::with_capacity(36);
    for (n, u, w) in faces {
        let corner = |su: f32, sw: f32| -> Vertex {
            let p = [
                0.5 * n[0] + 0.5 * su * u[0] + 0.5 * sw * w[0],
                0.5 * n[1] + 0.5 * su * u[1] + 0.5 * sw * w[1],
                0.5 * n[2] + 0.5 * su * u[2] + 0.5 * sw * w[2],
            ];
            v(p, n, [(su + 1.0) * 0.5, (sw + 1.0) * 0.5])
        };
        let (a, b, c, d) = (corner(-1.0, -1.0), corner(1.0, -1.0), corner(1.0, 1.0), corner(-1.0, 1.0));
        tris.extend_from_slice(&[a, b, c, a, c, d]);
    }

    let mut edges = Vec::with_capacity(24);
    let c = |x: f32, y: f32, z: f32| v([x * 0.5, y * 0.5, z * 0.5], [0.0; 3], [0.0; 2]);
    for &(a, b) in &[
        // bottom, top, verticals
        ((-1., -1., -1.), (1., -1., -1.)),
        ((1., -1., -1.), (1., -1., 1.)),
        ((1., -1., 1.), (-1., -1., 1.)),
        ((-1., -1., 1.), (-1., -1., -1.)),
        ((-1., 1., -1.), (1., 1., -1.)),
        ((1., 1., -1.), (1., 1., 1.)),
        ((1., 1., 1.), (-1., 1., 1.)),
        ((-1., 1., 1.), (-1., 1., -1.)),
        ((-1., -1., -1.), (-1., 1., -1.)),
        ((1., -1., -1.), (1., 1., -1.)),
        ((1., -1., 1.), (1., 1., 1.)),
        ((-1., -1., 1.), (-1., 1., 1.)),
    ] {
        edges.push(c(a.0, a.1, a.2));
        edges.push(c(b.0, b.1, b.2));
    }
    Mesh { tris, edges }
}

/// UV sphere, same parametrization as three.js SphereGeometry.
pub fn sphere() -> Mesh {
    let (ws, hs) = (SPHERE_WIDTH_SEGMENTS, SPHERE_HEIGHT_SEGMENTS);
    let pi = std::f32::consts::PI;

    let mut grid: Vec<Vec<Vertex>> = Vec::with_capacity(hs + 1);
    for iy in 0..=hs {
        let vv = iy as f32 / hs as f32;
        let mut row = Vec::with_capacity(ws + 1);
        for ix in 0..=ws {
            let u = ix as f32 / ws as f32;
            let p = [
                -SPHERE_RADIUS * (u * 2.0 * pi).cos() * (vv * pi).sin(),
                SPHERE_RADIUS * (vv * pi).cos(),
                SPHERE_RADIUS * (u * 2.0 * pi).sin() * (vv * pi).sin(),
            ];
            let n = [p[0] / SPHERE_RADIUS, p[1] / SPHERE_RADIUS, p[2] / SPHERE_RADIUS];
            row.push(v(p, n, [u, 1.0 - vv]));
        }
        grid.push(row);
    }

    let mut tris = Vec::new();
    for iy in 0..hs {
        for ix in 0..ws {
            let a = grid[iy][ix + 1];
            let b = grid[iy][ix];
            let c = grid[iy + 1][ix];
            let d = grid[iy + 1][ix + 1];
            if iy != 0 {
                tris.extend_from_slice(&[a, b, d]);
            }
            if iy != hs - 1 {
                tris.extend_from_slice(&[b, c, d]);
            }
        }
    }

    // latitude/longitude lines
    let mut edges = Vec::new();
    for iy in 1..hs {
        for ix in 0..ws {
            edges.push(grid[iy][ix]);
            edges.push(grid[iy][ix + 1]);
        }
    }
    for ix in 0..ws {
        for iy in 0..hs {
            edges.push(grid[iy][ix]);
            edges.push(grid[iy + 1][ix]);
        }
    }
    Mesh { tris, edges }
}
