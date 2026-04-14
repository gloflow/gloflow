# gf_lang RPC System Reference

## Overview

The gf_lang RPC (Remote Procedure Call) system `enables distributed generative graphics across multiple nodes`. It's built directly into the language as system functions and uses an external API for transport.

---

## RPC Call

Invoke a function on a remote node.

### Syntax

```
["rpc_call", node_name, module_name, function_name, ["args", arg1, arg2, ...]]
```

### Parameters

- **node_name** (string): Target node identifier
- **module_name** (string): Module namespace on target node
- **function_name** (string): Function to invoke
- **args** (list): Arguments to pass to the function

### Return Value

Returns a map containing the function's result. 

### Example

```
// Simple call
["var", "$result", 
    ["rpc_call", "worker_1", "geometry", "generate", ["args", 10, 5.0]]
]

// Use result in expression
[["x", ["$result", "offset"]], "cube"]

// Distributed rendering
["rule", "distribute_work", ["iters_max", 4], [
    ["rpc_call", 
        ["concat", "worker_", "$i"],  // worker_0, worker_1, worker_2, worker_3
        "render",
        "render_section",
        ["args", "$i", 100]
    ],
    "distribute_work"
]]
```

---

## RPC Serve

Expose gf_lang functions as RPC endpoints.

### Syntax

```
["rpc_serve", node_name,
    ["handlers", [
        [url_path, module_name, function_name, ["args_spec", ...], ["code", [...]]],
        [url_path, module_name, function_name, ["args_spec", ...], ["code", [...]]],
        ...
    ]]
]
```

### Parameters

- **node_name** (string): Name of this server node
- **handlers** (list): List of handler definitions

### Handler Definition

Each handler is a 5-element list:

1. **url_path** (string): HTTP endpoint path (e.g., "/api/generate")
2. **module_name** (string): Module namespace
3. **function_name** (string): Function name
4. **args_spec** (list): Argument specification `["args_spec", "arg1", "arg2", ...]`
5. **code** (list): gf_lang AST to execute `["code", [...expressions...]]`

### Example

```
["rpc_serve", "geometry_server",
    ["handlers", [
        // Simple handler
        [
            "/generate/cube",
            "geometry",
            "create_cube",
            ["args_spec", "size"],
            ["code", [
                [["sx", "$size"], ["sy", "$size"], ["sz", "$size"], "cube"]
            ]]
        ],
        
        // Complex recursive handler
        [
            "/generate/tower",
            "geometry",
            "create_tower",
            ["args_spec", "height", "taper"],
            ["code", [
                ["rule", "tower", ["iters_max", "$height"], [
                    [
                        ["y", 1.2],
                        ["sx", ["*", "$taper", ["$i"]]],
                        ["sy", ["*", "$taper", ["$i"]]],
                        ["sz", ["*", "$taper", ["$i"]]],
                        "cube"
                    ],
                    "tower"
                ]],
                "tower"
            ]]
        ],
        
        // Parametric spiral
        [
            "/generate/spiral",
            "geometry",
            "create_spiral",
            ["args_spec", "count", "radius", "height"],
            ["code", [
                ["rule", "spiral", ["iters_max", "$count"], [
                    [
                        ["x", ["*", ["cos", ["*", "$i", 0.2]], "$radius"]],
                        ["z", ["*", ["sin", ["*", "$i", 0.2]], "$radius"]],
                        ["y", ["*", "$i", ["$height"]]],
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

---

## Implementation Types

### Client Function

```go
type GFrpcCallFun func(
    string,        // node name
    string,        // module name
    string,        // function name
    []interface{}  // arguments
) map[string]interface{}  // result map
```

### Server Function

```go
type GFrpcServeFun func(
    string,                // node name
    []*GFrpcServerHandler, // handlers
    GFexternAPI            // extern API
)
```

### Handler Structure

```go
type GFrpcServerHandler struct {
    URLpathStr  string        // HTTP endpoint path
    ModuleStr   string        // module namespace
    FunctionStr string        // function name
    ArgsSpecLst []interface{} // argument specification
    CodeASTlst  []interface{} // gf_lang code to execute
}
```

---

## Use Cases

### 1. Distributed Rendering Farm

```
// Master node distributes sections
["rule", "master", ["iters_max", 8], [
    ["var", "$section", ["rpc_call", 
        ["concat", "worker_", "$i"],
        "render",
        "render_section",
        ["args", "$i", ["*", "$i", 10], ["*", ["$i"], 10]]
    ]],
    "master"
]]
```

### 2. Microservices Architecture

```
// Geometry service
["rpc_serve", "geometry_service",
    ["handlers", [
        ["/api/primitive/cube", "prim", "cube", [...], [...]],
        ["/api/primitive/sphere", "prim", "sphere", [...], [...]],
        ["/api/complex/building", "complex", "building", [...], [...]]
    ]]
]

// Texture service
["rpc_serve", "texture_service",
    ["handlers", [
        ["/api/material/wood", "mat", "wood", [...], [...]],
        ["/api/material/metal", "mat", "metal", [...], [...]]
    ]]
]

// Client combines services
["var", "$geom", ["rpc_call", "geometry_service", "complex", "building", ["args", 100]]]
["var", "$mat", ["rpc_call", "texture_service", "mat", "wood", ["args"]]]
```

### 3. Procedural Content Pipeline

```
// Stage 1: Base generation
["var", "$base", ["rpc_call", "generator", "base", "create", ["args", 50]]]

// Stage 2: Transform
["var", "$transformed", ["rpc_call", "transformer", "ops", "spiral", ["args", "$base", 10]]]

// Stage 3: Detail pass
["var", "$detailed", ["rpc_call", "detailer", "add", "features", ["args", "$transformed"]]]

// Stage 4: Render
["rpc_call", "renderer", "render", "scene", ["args", "$detailed"]]
```

### 4. Real-time Generative API

```
// HTTP endpoints for web clients
["rpc_serve", "api_server",
    ["handlers", [
        ["/api/v1/generate/architecture", "arch", "gen", ["args_spec", "style", "size"], [...]],
        ["/api/v1/generate/nature", "nature", "gen", ["args_spec", "type", "complexity"], [...]],
        ["/api/v1/generate/abstract", "abstract", "gen", ["args_spec", "seed", "iterations"], [...]]
    ]]
]
```

---

## Network Architecture

### Star Topology

```
         Master Node
            /  |  \
           /   |   \
          /    |    \
    Worker1 Worker2 Worker3
```

Master distributes work via `rpc_call`, workers execute and return results.

### Pipeline Topology

```
Generator → Transformer → Detailer → Renderer
```

Each node calls next stage via `rpc_call`, passing results forward.

### Mesh Topology

```
    Node1 ←→ Node2
      ↕         ↕
    Node3 ←→ Node4
```

Nodes communicate peer-to-peer, all expose endpoints via `rpc_serve`.

---

## External Implementation Requirements

Host environments must provide:

### 1. Network Transport
- HTTP, WebSockets, gRPC, or custom protocol
- Connection management and pooling
- Error handling and retries

### 2. Service Discovery
- Static configuration (hostnames/IPs)
- Dynamic registry (Consul, etcd, etc.)
- DNS-based discovery

### 3. Request Routing
- Map incoming requests to handlers by URL path
- Extract module, function, and arguments
- Route to appropriate handler code

### 4. Code Execution
- Execute handler `CodeASTlst` using `gf_lang.Run()`
- Pass arguments as variables in execution context
- Capture and serialize results

### 5. Serialization
- Convert gf_lang types to wire format (JSON, MessagePack, Protocol Buffers)
- Handle nested structures (lists, maps)
- Preserve type information

### 6. Security (Optional)
- Authentication/authorization
- Rate limiting
- Input validation

---

## Best Practices

### 1. Error Handling

```
["var", "$result", ["rpc_call", "node", "module", "func", ["args", ...]]]
["if", ["==", "$result", "nil"],
    // Handle error - fallback behavior
    [["cr", 1.0], "cube"]  // red cube indicates error
]
```

### 2. Timeout Management

Implement timeouts in external API to prevent hung connections.

### 3. Result Caching

Cache frequently called RPC results to reduce network overhead.

### 4. Load Balancing

Distribute calls across multiple workers:

```
["var", "$worker", ["concat", "worker_", ["%", "$i", 4]]]  // Round-robin 4 workers
["rpc_call", "$worker", "module", "function", ["args", ...]]
```

### 5. Idempotency

Design handlers to be idempotent - same input always produces same output.

### 6. Logging

Log all RPC calls for debugging distributed systems.

---

## Debugging

### Enable Debug Output

The RPC system uses `spew.Dump()` for debugging handler loading:

```go
func loadHandlers(pHandlersExprsLst []interface{}) {
    spew.Dump(pHandlersExprsLst)  // Dumps handler structure
    // ...
}
```

### Common Issues

**Problem**: Handler not found
- **Solution**: Check URL path matches exactly, including leading `/`

**Problem**: Arguments not available in handler
- **Solution**: Verify `args_spec` names match variable references in code

**Problem**: Result not returned
- **Solution**: Ensure handler code returns a value (use `return` statement)

**Problem**: Serialization errors
- **Solution**: Check that result types are JSON-serializable

---

## Example: Complete Distributed System

```
// ========================================
// Node 1: Master Coordinator
// ========================================
[
    ["lang_v", "0.0.6"],
    
    // Distribute work to 4 workers
    ["rule", "distribute", ["iters_max", 4], [
        ["var", "$worker", ["concat", "worker_", "$i"]],
        ["var", "$result", ["rpc_call", 
            "$worker",
            "geometry",
            "generate_section", 
            ["args", "$i", 25]
        ]],
        "distribute"
    ]],
    
    "distribute"
]

// ========================================
// Nodes 2-5: Worker Nodes
// ========================================
[
    ["lang_v", "0.0.6"],
    
    ["rpc_serve", ["concat", "worker_", "$node_id"],
        ["handlers", [
            [
                "/generate/section",
                "geometry",
                "generate_section",
                ["args_spec", "section_id", "complexity"],
                ["code", [
                    ["rule", "section", ["iters_max", "$complexity"], [
                        [
                            ["x", ["*", "$section_id", 10]],
                            ["y", ["*", "$i", 0.5]],
                            ["cr", ["*", "$i", 0.05]],
                            "cube"
                        ],
                        "section"
                    ]],
                    "section"
                ]]
            ]
        ]]
    ]
]
```

---

## Performance Considerations

### Network Overhead
- RPC calls add network latency (typically 1-50ms)
- Batch operations when possible
- Use local computation for tight loops

### Serialization Cost
- Large data structures slow serialization
- Consider streaming for huge datasets
- Use binary formats (MessagePack) over JSON for speed

### Connection Pooling
- Reuse connections to avoid TCP handshake overhead
- Implement connection pools in external API

### Concurrency
- Make parallel RPC calls for independent operations
- External API should support concurrent requests

---

## Summary

The gf_lang RPC system provides:

- **Language-level integration**: RPC as first-class system functions
- **Flexible topology**: Support for any network architecture
- **Code mobility**: Ship executable gf_lang code to remote nodes
- **Clean separation**: Language defines protocol, host provides transport
- **Scalability**: Distribute heavy generative workloads across multiple machines

Ideal for distributed rendering farms, microservices architectures, and real-time procedural content generation at scale.
