[
    ["lang_v", "0.0.3"],
    ["set", "color-background", ["rgb", 0.9, 0.9, 0.9]],
    [["set", "material", ["wireframe", true]], ["set", "color", "#42b0f5"], ["sphere"]],
    [["set", "material", ["shader", "gf_shader_test"]], ["set", "material_prop", ["gf_shader_test", "shader_uniform", ["i", 0.3]]], ["x", 30], "sphere"],
    [["set", "material", ["shader", "gf_shader_test"]], ["set", "material_prop", ["gf_shader_test", "shader_uniform", ["i", 0.5]]], ["sx", -0.2, "sy", -0.2, "sz", -0.2], ["x", 40, "rz", 0.4], "sphere"],
    [["set", "material", ["shader", "gf_shader_test"]], ["set", "material_prop", ["gf_shader_test", "shader_uniform", ["i", 0.7]]], ["sx", -0.4, "sy", -0.4, "sz", -0.4], ["x", 50, "rz", 0.6], "sphere"],
    [["set", "material", ["shader", "gf_shader_test"]], ["set", "material_prop", ["gf_shader_test", "shader_uniform", ["i", ["*", 2, 0.35]]]], ["sx", 0.6, "sy", 0.6, "sz", 0.6], ["x", -20, "y", 20, "rx", 1], "sphere"],
    [["set", "material_prop", ["gf_shader_test", "shader_uniform", ["i", 0.9]]], "R1"],
    ["rule", "R1", [["iters_max", 5]], [["set", "color", "#42b0f5"], [["y", 10, "ry", 0.5], "cube"], ["R1"]]],
    ["shader", "gf_shader_test", ["uniforms", [["i", "float", 0.1]]], ["vertex", "varying vec3 fNormal;\nvarying vec3 worldPos;\nvarying vec3 localPos;\nvoid main() {\nfNormal  = normalize(normalMatrix * normal);\nvec4 pos = modelViewMatrix * vec4(position, 1.0);\nworldPos = pos.xyz;\nlocalPos = position;\ngl_Position = projectionMatrix * pos;\n}"], ["fragment", "precision highp float;\nuniform float i;\nvarying vec3 fNormal;\nvarying vec3 localPos;\nfloat pulse(float val, float dst) {\nreturn floor(mod(val * dst, 1.0) + i);\n}\nvoid main() {\nvec3 dir   = vec3(0.5, 1, 0);\nvec3 cpos  = localPos;\nvec3 color = vec3(1, pulse(cpos.y, 10.0), 1);\nfloat diffuse = .5 + dot(fNormal, dir);\ngl_FragColor  = vec4(diffuse * color, 1.0);\n}\n/*void main() {\ngl_FragColor = vec4(1, 1, 0, 1.0);\n}*/"]]
]
