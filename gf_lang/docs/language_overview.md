# gf_lang Language Overview

## Introduction

gf_lang is a domain-specific language (DSL) for generative 3D graphics, inspired by context-free grammars and procedural generation systems. Programs are written as recursive rule definitions that incrementally modify spatial state to generate complex 3D scenes from simple primitives.

## Language Philosophy

- **State-based execution**: All geometric properties (position, rotation, scale, color) are maintained in state that flows through rule evaluations
- **Recursive composition**: Complex forms emerge from recursive application of simple transformation rules
- **Declarative syntax**: Programs describe transformations rather than explicit control flow
- **Iteration control**: Built-in safeguards prevent infinite recursion through global and per-rule iteration limits

---

## Execution Model

### Tree-Based Interpretation

gf_lang programs are represented as Abstract Syntax Trees (ASTs) and interpreted recursively:

1. **Entry Point**: Execution begins with `executeTree()` receiving an expression AST (represented as a list)
2. **State Creation**: Each tree descent creates a new child `GFstate` that inherits parent values
3. **Sequential Evaluation**: Expression elements are evaluated left-to-right
4. **Recursive Descent**: Sub-expressions trigger recursive `executeTree()` calls
5. **State Merging**: Child states merge back into parents via `stateMergeChild()`

### State Structure

The `GFstate` struct contains all spatial and control information:

```go
type GFstate struct {
    // Spatial properties
    Xf, Yf, Zf                    float64  // position
    RotationXf, RotationYf, RotationZf float64  // rotation
    ScaleXf, ScaleYf, ScaleZf     float64  // scale
    ColorRedF, ColorGreenF, ColorBlueF float64  // color (0-1 range)
    
    // Iteration control
    ItersMaxInt           int      // global max iterations
    ItersNumGlobalInt     int      // total iterations executed
    RulesItersNumStackLst []int    // per-rule iteration counters
    
    // Execution state
    RulesNamesStackLst    []string           // call stack
    VarsMap               map[string]*GFvariableVal  // variables
    AnimationsActiveMap   map[string]interface{}     // active animations
}
```

### State Propagation

**Inheritance Pattern:**
- Every `executeTree()` call creates a new state via `stateCreateNew(pStateParent, ...)`
- New state **copies** all parent property values (position, rotation, scale, color, variables)
- Modifications in child scope don't affect parent until merge

**Merge Pattern:**
- After child expression completes, `stateMergeChild()` copies child state back to parent
- All spatial properties, variables, and iteration counters propagate upward
- Animations are NOT merged (only propagate downward)

**State Families:**
- `push coord_origin` creates a new coordinate system family
- Family maintains separate spatial context while inheriting iteration counters
- `pop coord_origin` restores previous family state

### Expression Evaluation Flow

```
executeTree(expressionLst, parentState, ...)
  │
  ├─ Create new child state from parent
  │
  ├─ For each element in expression:
  │   │
  │   ├─ Property modifier (x, y, rx, etc)?
  │   │   └─ Increment state property value
  │   │
  │   ├─ Sub-expression (nested list)?
  │   │   ├─ Recursive executeTree() call
  │   │   ├─ Get child state + optional return value
  │   │   └─ Merge child state into current state
  │   │
  │   ├─ Conditional (if)?
  │   │   ├─ Evaluate condition
  │   │   └─ Execute branch if true
  │   │
  │   ├─ State setter (set/push/pop)?
  │   │   └─ Modify state or external renderer state
  │   │
  │   ├─ Variable assignment?
  │   │   └─ Store in state.VarsMap
  │   │
  │   ├─ Rule call?
  │   │   ├─ System rule (cube/sphere/line)?
  │   │   │   └─ Call extern API to create geometry
  │   │   └─ User rule?
  │   │       ├─ Push to call stack
  │   │       ├─ Check iteration limits
  │   │       ├─ Execute rule body recursively
  │   │       └─ Pop from call stack
  │   │
  │   └─ Return statement?
  │       └─ Break execution, return value up call chain
  │
  └─ Return final state + optional value to parent
```

### Recursion and Iteration Limits

**Global Limit:**
- `ItersMaxInt` (default: 250) limits total rule invocations
- Prevents runaway recursion across all rules
- When reached, rules exit immediately

**Per-Rule Limit:**
- Rules can specify `iters_max` modifier: `rule <name> [iters_max: 50] { ... }`
- Tracks iterations within a specific rule's recursive calls
- Uses `RulesItersNumStackLst` stack to count per-rule depth

**Iteration Counter (`$i`):**
- System variable tracking current rule's iteration count
- Incremented on each rule invocation
- Scoped to each rule (doesn't propagate between different rules)
- Accessible in expressions for conditional logic

---

## 3D Primitives and Spatial System

### System Rules (Built-in Primitives)

Three fundamental geometric primitives:

1. **`cube`** - Creates a unit cube at current state position/rotation/scale
2. **`sphere`** - Creates a unit sphere at current state position/rotation/scale  
3. **`line`** - Creates a line segment at current state position/rotation/scale

**Execution:**
- System rules call external API functions (`CreateCubeFun`, `CreateSphereFun`, `CreateLineFun`)
- Pass all 12 spatial parameters: `(x, y, z, rx, ry, rz, sx, sy, sz, cr, cg, cb)`
- Interpreter delegates actual rendering to external engine (Three.js, Babylon.js, etc.)

### Spatial Properties

The language exposes 12 spatial properties modifiable through state:

#### Position
- `x` - X-axis position (left/right)
- `y` - Y-axis position (up/down)
- `z` - Z-axis position (forward/back)

#### Rotation
- `rx` - Rotation around X-axis (pitch)
- `ry` - Rotation around Y-axis (yaw)
- `rz` - Rotation around Z-axis (roll)

#### Scale
- `sx` - X-axis scale
- `sy` - Y-axis scale
- `sz` - Z-axis scale

#### Color
- `cr` - Red channel (0.0 - 1.0)
- `cg` - Green channel (0.0 - 1.0)
- `cb` - Blue channel (0.0 - 1.0)

### Property Modifiers

Properties are modified **incrementally** (additive):

```
[["x", 2.0], "cube"]          // move +2 on x-axis, create cube
[["x", 1.0], ["y", 1.0], "cube"]  // move +1 on x and y, create cube
```

**Implementation:**
- `statePropFloatIncrement(state, propertyName, incrementValue)`
- Uses reflection to modify state fields dynamically
- New value = current value + increment

### Layout Calculations

**Incremental Positioning:**
Since modifiers are additive, layouts compose naturally:

```
rule grid [iters_max: 9] {
    [["x", ["*", ["$i"], 2.0]], "cube"]  // x = $i * 2.0
    "grid"                               // recursive call
}
```
Creates cubes at x = 0, 2, 4, 6, 8, 10, 12, 14, 16

**Nested Transformations:**
State inheritance enables hierarchical layouts:

```
rule tower [iters_max: 5] {
    "cube"
    [["y", 1.1], ["sx", 0.9], ["sy", 0.9], "tower"]  // each level smaller and higher
}
```

### Coordinate Systems

**World Origin (Default):**
- All positions relative to global (0, 0, 0)
- State properties accumulate through execution

**Current Position Origin:**
```
["push", "coord_origin", "current_pos"]
    // ... operations using current position as origin
["pop", "coord_origin", "current_pos"]
```

- Creates new "state family" with independent coordinate system
- Spatial properties reset but iteration counters inherited
- Enables local coordinate transformations
- Stack-based: push/pop for nested contexts

### Spatial Transformations Example

```
rule spiral [iters_max: 50] {
    [
        ["x", ["cos", ["*", "$i", 0.2]]],      // x = cos(i * 0.2)
        ["z", ["sin", ["*", "$i", 0.2]]],      // z = sin(i * 0.2)
        ["y", ["*", "$i", 0.1]],               // y = i * 0.1
        ["ry", 0.1],                           // rotate around y-axis
        ["sx", 0.8], ["sy", 0.8], ["sz", 0.8], // scale down
        "cube"
    ]
    "spiral"
}
```

Creates ascending spiral of cubes using trigonometry and `$i` variable.

---

## Rule System

### User Rules

Rules are named expression templates that can be invoked recursively:

```
rule <name> [modifiers] {
    <expression1>
    <expression2>
    ...
}
```

**Multiple Definitions:**
- Rules can have multiple bodies: random selection at each invocation
- Enables stochastic/organic generation patterns

**Rule Invocation:**
1. Current rule name stored in call stack (`RulesNamesStackLst`)
2. New state created inheriting parent values
3. If entering new rule (not recursive), new iteration counter pushed
4. Random definition selected if multiple exist
5. Iteration limits checked (global and per-rule)
6. Rule body executed via recursive `executeTree()`
7. State merged back or replaced depending on recursion type

**Recursion Handling:**
- Same-rule recursion: continues same iteration counter
- Different-rule call: starts new iteration counter
- Stack restoration on limit reached or natural exit

### System Rules

Hardcoded primitives executed via extern API:
- `cube`, `sphere`, `line`
- No body definition in language
- Directly create geometry in rendering engine

---

## Variables and Expressions

### Variables

**Declaration & Assignment:**
```
["var", "$myvar", 10.0]      // declare with initial value
["=", "$myvar", 5.0]         // assign new value
```

**System Variables:**
- `$i` - Current rule's iteration counter (auto-managed)

**Storage:**
- Stored in `state.VarsMap`
- Propagate through state merging
- Scoped to state family

### Arithmetic Expressions

Supported operators: `+`, `-`, `*`, `/`, `%`

```
["+", 5.0, 3.0]              // 8.0
["*", "$i", 2.0]             // $i * 2
["+", ["*", 2, 3], 5]        // (2*3) + 5 = 11
```

**Evaluation:**
- Detected by `arithmeticEval()`
- Nested expressions evaluated recursively
- Can reference variables and system functions

### System Functions

Built-in mathematical functions:

- `["cos", value]` - Cosine
- `["sin", value]` - Sine
- `["tan", value]` - Tangent
- `["abs", value]` - Absolute value
- `["sqrt", value]` - Square root
- `["rand", min, max]` - Random float in range

Example:
```
["x", ["*", ["cos", "$i"], 5.0]]  // x = cos($i) * 5
```

### Conditionals

```
["if", [operator, operand1, operand2], expression]
```

Operators: `==`, `!=`, `<`, `>`, `<=`, `>=`

Example:
```
["if", ["<", "$i", 5],
    [["cr", 1.0], "cube"]  // red cubes for first 5 iterations
]
```

---

## State Setters

Global state operations that affect rendering or coordinate systems:

### Set Operations

**Color:**
```
["set", "color", ["rgb", 1.0, 0.0, 0.0]]    // RGB values
["set", "color", "#ff0000"]                  // hex color
["set", "color-background", ["rgb", 0.5, 0.5, 0.5]]
```

**Scale (uniform):**
```
["set", "scale", 2.0]  // sets sx=2, sy=2, sz=2
```

**Iteration Limit:**
```
["set", "iters_max", 100]  // override global iteration limit
```

**Material:**
```
["set", "material", ["wireframe", true]]
["set", "material", ["shader", "customShader"]]
```

**Material Properties:**
```
["set", "material_prop", ["shaderName", "shader_uniform", "uniformName", value]]
```

### Coordinate System Control

**Push/Pop Origin:**
```
["push", "coord_origin", "current_pos"]
    // new coordinate system at current position
    [["x", 5], "cube"]  // x=5 relative to pushed origin
["pop", "coord_origin", "current_pos"]
```

Creates isolated spatial contexts for complex hierarchical structures.

---

## Animation System

Animate properties over time:

```
["animate", 
    [
        ["property_to_animate", "x"],
        ["start_value", 0.0],
        ["end_value", 10.0]
    ],
    ["duration_sec", 2.0],
    ["repeat", true]
]
```

**Behavior:**
- Calls `pExternAPI.AnimateFun()` 
- Delegates to external animation system
- Properties animated: any state property (`x`, `y`, `z`, `rx`, `ry`, `rz`, `sx`, `sy`, `sz`, `cr`, `cg`, `cb`)

---

## Program Structure

### Basic Program

```
[
    ["lang_v", "0.0.6"],
    
    ["rule", "main", [
        ["x", 1.0],
        "cube"
    ]],
    
    "main"
]
```

### Complex Program Example

```
[
    ["lang_v", "0.0.6"],
    
    // Background color
    ["set", "color-background", ["rgb", 0.1, 0.1, 0.15]],
    
    // Define recursive tower rule
    ["rule", "tower", ["iters_max", 10], [
        [
            ["cr", ["*", "$i", 0.1]],  // progressively more red
            ["y", 1.2],                 // move up
            ["sx", 0.9], ["sy", 0.9], ["sz", 0.9],  // scale down
            "cube"
        ],
        "tower"  // recurse
    ]],
    
    // Define spiral of towers
    ["rule", "spiral", ["iters_max", 8], [
        ["push", "coord_origin", "current_pos"],
            [
                ["x", ["*", ["cos", ["*", "$i", 0.785]], 5]],
                ["z", ["*", ["sin", ["*", "$i", 0.785]], 5]],
                "tower"
            ],
        ["pop", "coord_origin", "current_pos"],
        "spiral"
    ]],
    
    // Start execution
    "spiral"
]
```

---

## RPC System

gf_lang includes a built-in RPC (Remote Procedure Call) system for distributed generative graphics and inter-node communication.

### Architecture

The RPC system operates in two modes:

1. **RPC Call** - Invoke functions on remote nodes
2. **RPC Serve** - Expose functions as RPC endpoints

Both are integrated directly into the language and delegate to external implementations via the `GFexternAPI`.

### RPC Call

Call a function on a remote node from within gf_lang code:

```
["rpc_call", 
    "node_name",      // target node identifier
    "module_name",    // module containing function
    "function_name",  // function to invoke
    ["args", arg1, arg2, ...]  // argument list
]
```

**Example:**
```
["var", "$result", 
    ["rpc_call", 
        "worker_node_1",
        "geometry",
        "generate_spiral",
        ["args", 10, 2.0, 0.5]
    ]
]
```

**Implementation:**
- `rpcCallEval()` extracts node, module, function, and arguments
- Delegates to `pExternAPI.RPCcall()`
- Returns result map that can be assigned to variables
- Host environment handles actual network communication

**Function Signature:**
```go
type GFrpcCallFun func(
    string,        // node name
    string,        // module name
    string,        // function name
    []interface{}  // arguments
) map[string]interface{}  // result
```

### RPC Serve

Expose gf_lang functions as RPC endpoints that can be called by other nodes:

```
["rpc_serve", 
    "node_name",
    ["handlers", [
        [
            "/path/to/endpoint",
            "module_name",
            "function_name",
            ["args_spec", ...],
            ["code", [
                // gf_lang code to execute
            ]]
        ],
        // ... more handlers
    ]]
]
```

**Example:**
```
["rpc_serve", "geometry_server",
    ["handlers", [
        [
            "/generate/tower",
            "geometry",
            "create_tower",
            ["args_spec", "height", "width"],
            ["code", [
                ["rule", "tower", ["iters_max", ["$height"]], [
                    [["y", 1.2], ["sx", 0.9], ["sy", 0.9], "cube"],
                    "tower"
                ]],
                "tower"
            ]]
        ],
        [
            "/generate/spiral",
            "geometry", 
            "create_spiral",
            ["args_spec", "count", "radius"],
            ["code", [
                ["rule", "spiral", ["iters_max", ["$count"]], [
                    [
                        ["x", ["*", ["cos", ["*", "$i", 0.2]], ["$radius"]]],
                        ["z", ["*", ["sin", ["*", "$i", 0.2]], ["$radius"]]],
                        "cube"
                    ],
                    "spiral"
                ]],
                "spiral"
            ]]
        ]
    ]]
]
```

### RPC Server Handler Structure

Each handler is defined as:

```go
type GFrpcServerHandler struct {
    URLpathStr  string        // HTTP endpoint path
    ModuleStr   string        // module namespace
    FunctionStr string        // function name
    ArgsSpecLst []interface{} // argument specification
    CodeASTlst  []interface{} // gf_lang AST to execute
}
```

**Handler Execution Flow:**
1. External RPC request arrives at host environment
2. Host routes to appropriate handler by URL path
3. Handler's `CodeASTlst` (gf_lang code) is executed
4. Arguments from RPC call are available as variables
5. Execution result returned to caller

### Loading Handlers

`loadHandlers()` parses handler definitions:

1. Validates handler structure (must be 5 elements)
2. Extracts URL path, module, function name
3. Parses code block (must start with `["code", ...]`)
4. Creates `GFrpcServerHandler` structs
5. Returns list of handlers to extern API

**Validation Rules:**
- Code block must be length 2: `["code", [...]]`
- First element must be string "code"
- Second element is the executable AST

### RPC Integration Pattern

**Client Side (RPC Call):**
```
gf_lang program → rpcCallEval() → GFexternAPI.RPCcall()
    → Network layer → Remote node
```

**Server Side (RPC Serve):**
```
Network request → Host routing → GFexternAPI.RPCserve()
    → Handler lookup → Execute CodeASTlst → Return result
```

### Use Cases

**Distributed Rendering:**
```
// Master node distributes work
["rule", "distribute", ["iters_max", 8], [
    ["rpc_call", 
        ["concat", "worker_", ["$i"]],  // worker_0, worker_1, etc.
        "render",
        "render_section",
        ["args", "$i", 100]
    ],
    "distribute"
]]
```

**Generative Pipeline:**
```
// Node 1: Generate base geometry
["var", "$base", ["rpc_call", "node1", "gen", "base", ["args", 10]]]

// Node 2: Apply transformations
["var", "$transformed", ["rpc_call", "node2", "transform", "spiral", ["args", "$base"]]]

// Node 3: Add details
["rpc_call", "node3", "detail", "add_features", ["args", "$transformed"]]
```

**Dynamic Content Server:**
```
// Serve procedurally generated 3D content via HTTP endpoints
["rpc_serve", "content_server",
    ["handlers", [
        ["/api/generate/architecture", "arch", "building", [...code...]],
        ["/api/generate/nature", "nature", "tree", [...code...]],
        ["/api/generate/abstract", "abstract", "pattern", [...code...]]
    ]]
]
```

### RPC Function Signatures

```go
// Call remote function
type GFrpcCallFun func(
    string,        // node identifier
    string,        // module name
    string,        // function name  
    []interface{}  // arguments list
) map[string]interface{}  // result map

// Serve RPC endpoints
type GFrpcServeFun func(
    string,                // node name
    []*GFrpcServerHandler, // handler definitions
    GFexternAPI            // API for handler execution
)
```

### External Implementation Requirements

Host environments must implement:

1. **Network Layer**: HTTP, WebSockets, gRPC, or custom protocol
2. **Node Discovery**: Service registry or static configuration
3. **Routing**: Map incoming requests to handlers by URL path
4. **Execution Context**: Run handler `CodeASTlst` with `gf_lang.Run()`
5. **Serialization**: Convert between gf_lang types and wire format

The language provides the syntax and structure; the host provides transport and infrastructure.

---

## Extern API Interface

The interpreter doesn't render directly - it calls external functions:

```go
type GFexternAPI struct {
    InitEngineFun   func(map[string]*GFshaderDef)
    SetStateFun     func(GFstateChange) []interface{}
    CreateCubeFun   GFgeometryFunc
    CreateSphereFun GFgeometryFunc
    CreateLineFun   GFgeometryFunc
    AnimateFun      func([]map[string]interface{}, float64, bool)
    RPCcall         GFrpcCallFun      // RPC client
    RPCserve        GFrpcServeFun     // RPC server
}
```

**Integration Pattern:**
1. Host environment (JavaScript + Three.js, Go + OpenGL, etc.) implements API
2. Passes API struct to `gf_lang.Run()`
3. Interpreter calls API functions with state parameters
4. Host creates actual 3D objects in rendering engine
5. RPC functions enable distributed execution across nodes

---

## Debugging Support</text>

<old_text line=520>
gf_lang is a specialized language for generative 3D graphics with:

1. **State-driven execution**: Spatial properties flow through recursive rule evaluations
2. **Incremental transformations**: Properties modified additively through expression evaluation  
3. **Coordinate system control**: Push/pop origins for hierarchical layouts
4. **Built-in iteration control**: Global and per-rule limits prevent infinite recursion
5. **External rendering**: Clean separation between language execution and graphics APIs
6. **Procedural generation**: Simple rules compose into complex 3D structures</text>
<new_text>
gf_lang is a specialized language for generative 3D graphics with:

1. **State-driven execution**: Spatial properties flow through recursive rule evaluations
2. **Incremental transformations**: Properties modified additively through expression evaluation  
3. **Coordinate system control**: Push/pop origins for hierarchical layouts
4. **Built-in iteration control**: Global and per-rule limits prevent infinite recursion
5. **External rendering**: Clean separation between language execution and graphics APIs
6. **Procedural generation**: Simple rules compose into complex 3D structures
7. **Built-in RPC system**: Distributed execution and inter-node communication for scalable generative pipelines</text>


Optional debug tracking via `GFprogramDebug`:

- **State History**: `StateHistoryLst` - all states created during execution
- **Entity Output**: `EntitiesOutputMap` - all primitives generated
- **Rule Calls**: `RulesCallsCounterMap` - invocation counts per rule
- **State Changes**: `ExternStateChangesLst` - all extern state modifications

Enables visualization of execution flow and generated geometry.

---

## Summary

gf_lang is a specialized language for generative 3D graphics with:

1. **State-driven execution**: Spatial properties flow through recursive rule evaluations
2. **Incremental transformations**: Properties modified additively through expression evaluation  
3. **Coordinate system control**: Push/pop origins for hierarchical layouts
4. **Built-in iteration control**: Global and per-rule limits prevent infinite recursion
5. **External rendering**: Clean separation between language execution and graphics APIs
6. **Procedural generation**: Simple rules compose into complex 3D structures

The execution model's recursive descent with state merging enables elegant expression of complex generative patterns while maintaining predictable control flow and resource limits.