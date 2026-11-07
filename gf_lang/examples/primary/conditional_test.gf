[
    ["lang_v", "0.0.3"],
    ["set", "color-background", ["rgb", 0.7, 0.7, 0.7]],
    [["set", "color", ["rgb", 0.7, 0, 0]], ["set", "material", ["wireframe", true]], ["set", "iters_max", 30], ["R1"]],
    ["rule", "R1", [["iters_max", 10]], [["print", ["first condition", "$i"]], ["if", ["==", "$i", 2], [["print", ["$i", "running R2 A"]], ["R2"]]], ["print", ["second condition", "$i"]], ["if", ["==", "$i", 4], [["print", ["$i", "running R2 B"]], ["R2"]]], ["print", ["third condition", "$i"]], ["if", ["==", "$i", 8], [["print", ["$i", "running R2 C"]], ["R2"]]], ["set", "color", ["rgb", 0.7, 0, 0]], ["*", 1, [["x", 2], "cube"]], ["R1"]]],
    ["rule", "R2", [["iters_max", 5]], [["print", ["===============in R2", "$i"]], ["y", 4], ["print", ["first R2 condition", "$i"]], ["if", [">", "$i", 2], [["set", "color", ["rgb", 0, 0.9, 0]], ["print", ["$i", "running R3"]], ["R3"]]], [["print", ["===============in R2 cube", "$i"]], ["set", "color", ["rgb", 0, 0, 0.9]], ["*", 1, ["cube"]], ["print", ["===============exit R2", "$i"]], ["R2"]]]],
    ["rule", "R3", [["iters_max", 5]], [["print", ["$i", "===============in R3"]], ["*", 1, [["z", 4], "cube"]], ["print", ["$i", "===============exit R3"]], "R3"]]
]
