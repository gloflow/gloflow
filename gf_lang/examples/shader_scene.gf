[
    ["lang_v", 0.3],
    ["set", "color-background", ["rgb", 0.9, 0.9, 0.9]],

    [
        ["set", "material", ["wireframe", true]],
        ["set", "color", ["rgb", 0.2, 0.3, 0.8]],
        ["sphere"]
    ],

    [
        ["set", "material", ["shader", "gf_shader_test"]],
        ["set", "material_prop", ["gf_shader_test", "shader_uniform", ["i", 0.3]]],
        ["x", 30], "sphere"
    ],
    [
        ["set", "material", ["shader", "gf_shader_test"]],
        ["set", "material_prop", ["gf_shader_test", "shader_uniform", ["i", 0.7]]],
        ["sx", -0.4, "sy", -0.4, "sz", -0.4],
        ["x", 50, "rz", 0.6], "sphere"
    ],
    [
        ["set", "material", ["shader", "gf_shader_test"]],
        ["set", "material_prop", ["gf_shader_test", "shader_uniform", ["i", 0.5]]],
        ["x", -25, "y", 15, "rx", 1.0], "cube"
    ],

    ["shader", "gf_shader_test",
        ["uniforms", [
            ["i", "float", 0.1]
        ]],
        ["vertex", "varying vec3 fNormal;\nvarying vec3 worldPos;\nvarying vec3 localPos;\nvoid main() {\nfNormal  = normalize(normalMatrix * normal);\nvec4 pos = modelViewMatrix * vec4(position, 1.0);\nworldPos = pos.xyz;\nlocalPos = position;\ngl_Position = projectionMatrix * pos;\n}"],
        ["fragment", "precision highp float;\nuniform float i;\nvarying vec3 fNormal;\nvarying vec3 localPos;\nfloat pulse(float val, float dst) {\nreturn floor(mod(val * dst, 1.0) + i);\n}\nvoid main() {\nvec3 dir   = vec3(0.5, 1, 0);\nvec3 cpos  = localPos;\nvec3 color = vec3(1, pulse(cpos.y, 10.0), 1);\nfloat diffuse = .5 + dot(fNormal, dir);\ngl_FragColor  = vec4(diffuse * color, 1.0);\n}"]
    ]
]
