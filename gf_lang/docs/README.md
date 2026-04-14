# gf_lang Documentation

## Overview

gf_lang is a domain-specific language (DSL) for generative 3D graphics and distributed procedural content generation. It combines recursive rule-based composition with spatial transformations, built-in iteration control, and native RPC capabilities for scalable distributed execution.

## Documentation Files

### [language_overview.md](language_overview.md)
**Comprehensive language reference covering:**
- Execution model and state propagation
- 3D primitives and spatial system
- Rule system (user and system rules)
- Variables, expressions, and conditionals
- State setters and coordinate systems
- Animation system
- RPC system architecture
- External API interface
- Debugging support

### [rpc_reference.md](rpc_reference.md)
**RPC system quick reference:**
- `rpc_call` - Client-side remote function invocation
- `rpc_serve` - Server-side endpoint exposure
- Handler definition structure
- Use cases and network topologies
- Implementation requirements
- Best practices and debugging

### [execution.md](execution.md)
**Detailed execution internals:**
- `executeTree` function flow
- Sub-expression evaluation and substitution
- Property modifiers and arithmetic
- Conditionals, state setters, and animations
- Rule dispatch (system vs user)
- State management and merging
- Iteration limits and recursion control

## Quick Start

### Basic Program Structure

```
[
    ["lang_v", "0.0.6"],
    
    // Define rules
    ["rule", "my_rule", [
        // transformations
        [["x", 1.0], "cube"]
    ]],
    
    // Execute
    "my_rule"
]
```

### Key Concepts

**1. State-Driven Execution**
- All spatial properties (position, rotation, scale, color) stored in state
- State flows through recursive rule evaluations
- Child states inherit parent values and merge back after execution

**2. Incremental Transformations**
- Properties modified additively: `["x", 2.0]` adds 2.0 to current x
- Transformations compose naturally through nesting
- Enables declarative layout descriptions

**3. Recursive Rules**
- Rules can call themselves or other rules
- Global and per-rule iteration limits prevent infinite recursion
- `$i` variable tracks current rule iteration

**4. Built-in Primitives**
- `cube` - Unit cube at current state
- `sphere` - Unit sphere at current state
- `line` - Line segment at current state

**5. RPC System**
- `rpc_call` - Invoke functions on remote nodes
- `rpc_serve` - Expose handlers as RPC endpoints
- Native distributed execution support

## Example Programs

### Simple Tower

```
[
    ["lang_v", "0.0.6"],
    
    ["rule", "tower", ["iters_max", 10], [
        [
            ["y", 1.2],                              // move up
            ["sx", 0.9], ["sy", 0.9], ["sz", 0.9],  // scale down
            ["cr", ["*", "$i", 0.1]],                // progressively red
            "cube"
        ],
        "tower"  // recurse
    ]],
    
    "tower"
]
```

### Parametric Spiral

```
[
    ["lang_v", "0.0.6"],
    
    ["rule", "spiral", ["iters_max", 50], [
        [
            ["x", ["*", ["cos", ["*", "$i", 0.2]], 5.0]],
            ["z", ["*", ["sin", ["*", "$i", 0.2]], 5.0]],
            ["y", ["*", "$i", 0.1]],
            ["ry", 0.1],
            "cube"
        ],
        "spiral"
    ]],
    
    "spiral"
]
```

### Distributed Rendering

```
[
    ["lang_v", "0.0.6"],
    
    // Master node distributes work
    ["rule", "distribute", ["iters_max", 4], [
        ["rpc_call",
            ["concat", "worker_", "$i"],
            "render",
            "render_section",
            ["args", "$i", 100]
        ],
        "distribute"
    ]],
    
    "distribute"
]
```

### RPC Server

```
[
    ["lang_v", "0.0.6"],
    
    ["rpc_serve", "geometry_server",
        ["handlers", [
            [
                "/generate/tower",
                "geometry",
                "create_tower",
                ["args_spec", "height"],
                ["code", [
                    ["rule", "tower", ["iters_max", "$height"], [
                        [["y", 1.2], "cube"],
                        "tower"
                    ]],
                    "tower"
                ]]
            ]
        ]]
    ]
]
```

## Spatial Properties

### Position
- `x` - X-axis (left/right)
- `y` - Y-axis (up/down)
- `z` - Z-axis (forward/back)

### Rotation
- `rx` - X-axis rotation (pitch)
- `ry` - Y-axis rotation (yaw)
- `rz` - Z-axis rotation (roll)

### Scale
- `sx` - X-axis scale
- `sy` - Y-axis scale
- `sz` - Z-axis scale

### Color
- `cr` - Red channel (0.0-1.0)
- `cg` - Green channel (0.0-1.0)
- `cb` - Blue channel (0.0-1.0)

## System Functions

### Mathematical
- `["cos", value]` - Cosine
- `["sin", value]` - Sine
- `["tan", value]` - Tangent
- `["abs", value]` - Absolute value
- `["sqrt", value]` - Square root
- `["rand", min, max]` - Random float

### Data Structures
- `["make", "list", ...]` - Create list
- `["make", "map", ...]` - Create map
- `["len", collection]` - Get length

### RPC
- `["rpc_call", node, module, function, args]` - Call remote function
- `["rpc_serve", node, handlers]` - Serve RPC endpoints

## Arithmetic Operators

- `["+", a, b]` - Addition
- `["-", a, b]` - Subtraction
- `["*", a, b]` - Multiplication
- `["/", a, b]` - Division
- `["%", a, b]` - Modulo

## Logic Operators

- `["==", a, b]` - Equal
- `["!=", a, b]` - Not equal
- `["<", a, b]` - Less than
- `[">", a, b]` - Greater than
- `["<=", a, b]` - Less than or equal
- `[">=", a, b]` - Greater than or equal

## State Setters

### Color
```
["set", "color", ["rgb", r, g, b]]
["set", "color", "#ff0000"]
["set", "color-background", ["rgb", r, g, b]]
```

### Scale
```
["set", "scale", 2.0]  // uniform scale
```

### Material
```
["set", "material", ["wireframe", true]]
["set", "material", ["shader", "shader_name"]]
```

### Coordinate System
```
["push", "coord_origin", "current_pos"]
    // operations in local coordinate system
["pop", "coord_origin", "current_pos"]
```

### Iteration Control
```
["set", "iters_max", 100]
```

## Animation

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

## Variables

### Declaration & Assignment
```
["var", "$myvar", 10.0]    // declare
["=", "$myvar", 5.0]       // assign
```

### System Variables
- `$i` - Current rule iteration counter

## Conditionals

```
["if", [operator, operand1, operand2],
    expression_if_true
]
```

Example:
```
["if", ["<", "$i", 5],
    [["cr", 1.0], "cube"]  // red cubes for i < 5
]
```

## Architecture

### Execution Flow

```
Program AST → executeTree()
    ↓
Create child state from parent
    ↓
For each expression element:
    - Property modifier → increment state
    - Sub-expression → recursive executeTree()
    - Rule call → execute rule body
    - State setter → modify external state
    - Variable → store in state
    - Conditional → evaluate and branch
    ↓
Merge child state to parent
    ↓
Return state + optional value
```

### External API Integration

```
Host Environment (Three.js, OpenGL, etc.)
    ↓
Implements GFexternAPI
    - InitEngineFun
    - SetStateFun
    - CreateCubeFun
    - CreateSphereFun
    - CreateLineFun
    - AnimateFun
    - RPCcall
    - RPCserve
    ↓
Pass to gf_lang.Run()
    ↓
Interpreter calls API functions
    ↓
Host creates actual 3D objects
```

## Use Cases

### Procedural Architecture
Generate buildings, cities, and structures through recursive rules.

### Organic Forms
Create natural patterns (trees, plants, corals) using stochastic rules.

### Abstract Art
Generate complex abstract compositions from simple primitives.

### Distributed Rendering
Split heavy generative workloads across multiple machines via RPC.

### Real-time Content API
Expose procedural generation as HTTP endpoints for dynamic content.

### Generative Pipelines
Chain multiple processing stages across specialized nodes.

## Implementation

### Language: Go
- Interpreter: `gloflow/gf_lang/go/gf_lang/`
- Type-safe AST processing
- Reflection-based state property access
- Clean separation between language and rendering

### Targets
- **WebAssembly**: Browser-based execution with Three.js
- **Native**: Server-side with OpenGL/Vulkan
- **Distributed**: Multi-node clusters with RPC

## Performance Characteristics

### Execution
- Tree-walking interpreter (not compiled)
- State copying on each descent (isolated scopes)
- Recursive function calls (bounded by iteration limits)

### Optimization Opportunities
- AST pre-expansion (constant folding, loop unrolling)
- State object pooling
- Memoization of pure function results
- Parallel RPC call execution

## Debugging

Enable debug mode to capture:
- **State History**: All states created during execution
- **Entity Output**: All primitives generated with properties
- **Rule Calls**: Invocation counts per rule
- **State Changes**: All external state modifications

## Language Philosophy

**Declarative over Imperative**
Describe transformations, not control flow.

**Composition over Configuration**
Complex forms emerge from simple rule composition.

**State over Objects**
Properties flow through evaluation, not attached to objects.

**Bounds over Freedom**
Built-in limits prevent runaway recursion.

**Distribution over Centralization**
Native RPC enables scalable architectures.

## Contributing

See main project repository for contribution guidelines.

## License

GNU General Public License v2.0 or later

---

**Version**: 0.0.6  
**Author**: Ivan Trajkovic  
**Project**: GloFlow

For detailed technical information, see individual documentation files linked above.