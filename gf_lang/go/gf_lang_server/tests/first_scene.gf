[
[
    ["lang_v", 0.3],
    ["set", "color-background", ["rgb", 0.5, 0.5, 0.5]],
    ["set", "material", ["wireframe", true]],

    [["set", "color", ["rgb", 0, 0, 1.0]], "cube"],                                   // cube
    [["set", "color", ["rgb", 0, 0, 1.0]], ["x", 2.0], "cube"],                       // {x 2} cube
    [["set", "color", ["rgb", 0, 0, 1.0]], ["x", 4.0], "cube"],                       // {x 4} cube
    [["set", "color", ["rgb", 0, 0, 1.0]], ["x", 6.0], "cube"],                       // {x 6} cube
    [["set", "color", ["rgb", 0, 0, 1.0]], ["x", 9.0, "z", 2.0], ["y", 3.0], "cube"], // {x 7 z 1} {y 3} cube
    
    [
        ["set", "color", ["rgb", 0, 0, 1.0]],
        ["*", 3, [["y", -2.0, "z", -2.0, "x", -2.0], "cube"]]  // 3 * {x -2 z -2 x -2} cube
    ],
    
    [
        ["set", "color", ["rgb", 0, 0, 1.0]],
        ["*", 3, [["x", 5.0], ["*", 6, [["y", 2.0], "cube"]]]] // 3 * {x 5} 6 * {y 2} cube
    ],

    
    // transforms to:
    // [
    // 
    // ["x", 5.0], [
    //     [["y", 2.0], "cube"]
    //     [["y", 2.0], "cube"]
    //     [["y", 2.0], "cube"]
    //     [["y", 2.0], "cube"]
    //     [["y", 2.0], "cube"]
    //     [["y", 2.0], "cube"]
    // ]
    // 
    // ["x", 5.0], [
    //     [["y", 2.0], "cube"]
    //     [["y", 2.0], "cube"]
    //     [["y", 2.0], "cube"]
    //     [["y", 2.0], "cube"]
    //     [["y", 2.0], "cube"]
    //     [["y", 2.0], "cube"]
    // ]
    // 
    // ["x", 5.0], [
    //     [["y", 2.0], "cube"]
    //     [["y", 2.0], "cube"]
    //     [["y", 2.0], "cube"]
    //     [["y", 2.0], "cube"]
    //     [["y", 2.0], "cube"]
    //     [["y", 2.0], "cube"]
    // ]
    // ]
    

    [
        ["set", "color", ["rgb", 0, 0, 1.0]],
        ["*", 3, [["x", 5.4], ["*", 1, [["z", 1.6], "cube"]]]] // 3 * {x 5} 6 * {y 2} cube
    ],
    [
        ["set", "color", ["rgb", 0, 0, 1.0]],
        ["*", 2, [["x", 2.1], ["*", 2, [["y", 2.1], "cube"]]]] // 3 * {x 5} 6 * {y 2} cube
    ],
    [
        ["set", "color", ["rgb", 0, 0, 1.0]],
        ["*", 6, [["x", 4.4], ["*", 1, [["z", 2.4], "cube"]]]] // 3 * {x 5} 6 * {y 2} cube
    ],
    [
        ["set", "color", ["rgb", 0, 0, 1.0]],
        ["*", 4, [["x", 5.2], ["*", 1, [["y", 2.2], "cube"]]]] // 3 * {x 5} 6 * {y 2} cube
    ],
    [
        ["set", "color", ["rgb", 0, 0, 1.0]],
        ["*", 1, [["x", 6.1], ["*", 3, [["z", 2.0], "cube"]]]] // 3 * {x 5} 6 * {y 2} cube
    ],
    [
        ["set", "color", ["rgb", 0, 0, 1.0]],
        ["*", 1, [["x", 1.4], ["*", 2, [["y", 2.1], "cube"]]]] // 3 * {x 5} 6 * {y 2} cube
    ]
]
]