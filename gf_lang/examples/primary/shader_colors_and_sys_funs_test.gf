[
    ["lang_v", "0.0.4"],
    ["set", "color-background", ["rgb", 0.9, 0.9, 0.9]],
    [["set", "color", ["rgb", 0.2, 0.2, 0.2]], ["set", "material", ["shader", "gf_shader_test"]], ["set", "material_prop", ["gf_shader_test", "shader_uniform", ["i", 0.8]]], ["set", "material_prop", ["gf_shader_test", "shader_uniform", ["cr", "$cr"]]], ["set", "material_prop", ["gf_shader_test", "shader_uniform", ["cg", "$cg"]]], ["set", "material_prop", ["gf_shader_test", "shader_uniform", ["cb", "$cb"]]], ["y", 2, "z", 2], "cube"],
    [["set", "material", ["shader", "gf_shader_test"]], ["set", "color", ["rgb", 0.2, 0.8, 0.2]], ["x", 2], "seed"],
    [["set", "material", ["shader", "gf_shader_test"]], ["set", "color", ["rgb", 0.5, 0.8, 0.2]], ["x", 4], "seed"],
    [["set", "material", ["shader", "gf_shader_test"]], ["set", "color", ["rgb", 0.5, 0.8, 0.6]], ["x", 6], "seed"],
    [["x", 8, "y", 4], ["set", "line", ["start"]], ["set", "scale", [0.3, 0.3, 0.3]], "line_draw"],
    [["x", 8.5, "y", 4], ["set", "line", ["start"]], ["set", "scale", [0.3, 0.3, 0.3]], "line_draw"],
    [["x", 9, "y", 4], ["set", "line", ["start"]], ["set", "scale", [0.3, 0.3, 0.3]], "line_draw"],
    [["x", 9.5, "y", 4], ["set", "line", ["start"]], ["set", "scale", [0.3, 0.3, 0.3]], "line_draw"],
    [["x", 10, "y", 4], ["set", "line", ["start"]], ["set", "scale", [0.3, 0.3, 0.3]], "line_draw"],
    [["x", 10.5, "y", 4], ["set", "line", ["start"]], ["set", "scale", [0.3, 0.3, 0.3]], "line_draw"],
    [["x", 11, "y", 4], ["set", "line", ["start"]], ["set", "scale", [0.3, 0.3, 0.3]], "line_draw"],
    ["rule", "line_draw", [["iters_max", 30]], [["if", ["==", "$i", 20], [["set", "line", ["start"]]]], [["z", ["rand", [0, 2]], "x", ["rand", [-1, 1]], "y", ["rand", [-0.7, 0.7]], "rz", 0.3], "line"], ["cube"], ["line_draw"]]],
    ["rule", "seed", [["iters_max", 10]], [[["set", "material_prop", ["gf_shader_test", "shader_uniform", ["cr", "$cr"]]], ["set", "material_prop", ["gf_shader_test", "shader_uniform", ["cg", "$cg"]]], ["set", "material_prop", ["gf_shader_test", "shader_uniform", ["cb", "$cb"]]], ["z", ["+", 1, ["rand", [0, 0.5]]], "cg", -0.05], "cube"], ["seed"]]],
    ["shader", "gf_shader_test", ["uniforms", [["i", "float", 0.1], ["cr", "float", 0.5], ["cg", "float", 0.5], ["cb", "float", 0.5]]], ["vertex", "varying vec3 normal_f;\nvarying vec3 local_pos_v3;\nvoid main() {\nnormal_f     = normalize(normalMatrix * normal);\nlocal_pos_v3 = position;\nvec4 pos     = modelViewMatrix * vec4(position, 1.0);\ngl_Position  = projectionMatrix * pos;\n}"], ["fragment", "precision highp float;\nuniform float i;\nuniform float cr;\nuniform float cg;\nuniform float cb;\nvarying vec3 normal_f;\nvarying vec3 local_pos_v3;\nvoid main() {\nvec3 light_direction_v3 = vec3(0.8, 1, 0.8);\nvec3 color_v3           = vec3(cr, cg, cb);\nfloat diffuse_f = .5 + dot(normal_f, light_direction_v3);\ngl_FragColor    = vec4(diffuse_f * color_v3, 1.0);\n}"]]
]
