["lang_v", 0.1],
["cube"], // cube

["R1"],
[["x", 2.0], "R1"],
[["x", 4.0,  "y", 1.0], "R1"],
[["x", 6.0,  "y", 2.0], "R1"],
[["x", 8.0,  "y", 3.0], "R1"],
[["x", 10.0, "y", 4.0], "R1"],
[["x", 12.0, "y", 5.0], "R1"],
[["x", 14.0, "y", 6.0], "R1"],
[["x", 16.0, "y", 7.0], "R1"],
[["x", 18.0, "y", 8.0], "R1"],
[["x", 20.0, "y", 9.0], "R1"],
[["x", 21.0, "y", 10.0], "R1"],
[["x", 22.0, "y", 11.0], "R1"],
[["x", 23.0, "y", 12.0], "R1"],
[["x", 24.0, "y", 13.0], "R1"],
["rule", "R1", [
    ["*", 2.0, [["x", 0.5, "y", 0.7, "rz", 0.07], "cube"]],
    ["R1"]
]],
["rule", "R1", [
    ["*", 1.0, [["x", 0.5, "y", 0.7, "z", 1.0, "rz", 0.05], "cube"]],
    ["R1"]
]],
["rule", "R1", [
    ["*", 1.0, [["x", 0.5, "y", 0.7, "z", 1.3, "ry", 0.05, "rz", 0.05], "cube"]],
    ["R1"]
]],
["rule", "R1", [
    ["*", 1.0, [["x", 0.5, "y", 0.7, "z", 1.5, "rx", 0.05, "rz", 0.05], "cube"]],
    ["R1"]
]]