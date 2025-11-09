Summary of how gf_lang is executed (based on `gf_lang/go/gf_lang/gf_lang_exec.go`)

High-level flow
- Execution starts by calling `executeTree` with an expression AST (a list), a parent state and environment (rules, shaders, extern API, debug).
- `executeTree` creates a new child `state` for that expression-tree descent and evaluates the expression elements sequentially. It uses recursion to evaluate nested sub-expressions.
- There are distinct expression kinds handled in the main loop: property modifiers, sub-expressions, arithmetic expressions, conditionals (`if`), state setters (`set`, `push`, `pop`), `print`, `animate`, variable assignment, `return`, and rule calls (user or system).
- Child subtree evaluations return both a child state and optionally a value. Child state is merged back into the parent via `stateMergeChild`.

Key entry points (examples)
```gf_lang/go/gf_lang/gf_lang_exec.go#L32-80
func executeTree(pExpressionASTlst []interface{},
    pStateParent         *GFstate,
    pRulesDefsMap        GFruleDefs,
    pShaderDefsMap       map[string]*GFshaderDef,
    pStateFamilyStackLst []*GFstate,
    pExternAPI           GFexternAPI,
    pDebug               *GFprogramDebug) (*GFstate, interface{}, error) {

    symbols := getSymbolsAndConstants()

    //--------------------
    // STATE_NEW
    // IMPORTANT!! - on every tree descent a new independent state is constructed
    state := stateCreateNew(pStateParent, pDebug)
```

Sub-expressions, substitution, and returns
- A sub-expression is evaluated by recursively calling `executeTree`. If the sub-expression returns a value (non-nil), `executeTree` substitutes that value into the parent expression list in place of the sub-expression and continues evaluation.
- If a sub-expression is a `return`, `executeTree` immediately returns the current state plus the returned value up the call stack.

Property modifiers and arithmetic
- Property modifiers (predefined properties) are incremented incrementally on the current state via `statePropFloatIncrement`. Modifiers can be direct numbers or computed by evaluating a sub-expression.
- If an expression is detected to be an arithmetic expression, `arithmeticEval` is used and the arithmetic result is returned as the expression value.

Conditionals, prints and animations
- `exprConditional` evaluates a 3-element logic expression (operator, operand1, operand2) and, if true, recursively executes the conditional branch via `executeTree` and then merges the child state back into the caller.
- `exprPrint` formats variable or literal strings and prints using `fmt.Printf`.
- `exprAnimation` computes start/end values for properties and calls `pExternAPI.AnimateFun` (and can set repeat).

Rules (user vs system)
- `exprRuleCall` handles invocation of rules:
  - System rules (in a predefined list) are handled by `exprRuleSysCall` and map directly to extern API primitives (e.g., `cube`, `sphere`, `line`) — these call `pExternAPI.CreateCubeFun(...)` and add the entity to the output.
  - User rules (defined by the program) create a new state that inherits the parent's values (a new state per rule invocation). The interpreter:
    - pushes the rule name onto a call stack,
    - increments iteration counters (global and per-rule),
    - picks a rule definition (random def if rule has multiple defs),
    - enforces iteration limits (global `iters_max` and rule-specific `iters_max`), and
    - recursively runs the rule body via `executeTree`.
  - When returning from a rule, it restores iteration stacks, pops rule names, and either returns the new state or merges child state depending on whether the call was to a different rule or recursive within the same rule.

Example of rule dispatch (system vs user)
```gf_lang/go/gf_lang/gf_lang_exec.go#L378-440
func exprRuleCall(pCalledRuleNameStr string,
    pExpressionLst       []interface{},
    pStateParent         *GFstate,
    pRulesDefsMap        GFruleDefs,
    pShaderDefsMap       map[string]*GFshaderDef,
    pStateFamilyStackLst []*GFstate,
    pExternAPI           GFexternAPI,
    pDebug               *GFprogramDebug) (*GFstate, error) {
...
    if gf_core.ListContainsStr(pCalledRuleNameStr, symbols.SystemRulesLst) {
        newState := exprRuleSysCall(pCalledRuleNameStr,
            pStateParent,
            pExternAPI,
            pDebug)
        return newState, nil
    } else if gf_core.MapHasKey(pRulesDefsMap, pCalledRuleNameStr) {
        // user rule: create new state, push call stack, pick def, execute
```

State management and merging
- Every descent (executeTree call or rule invocation) creates a new `GFstate` via `stateCreateNew`, inheriting values from its parent. This keeps states isolated across rule invocations.
- After a child subtree completes, its state is merged back into the parent using `stateMergeChild` — except in some rule exit cases where the state must be reset to the caller’s state.
- Some state setters (e.g., `push`/`pop` on `coord_origin`) can return a new state (special handling).

Evaluation helpers & extern API
- Arithmetic is handled by `arithmeticEval`. System functions are recognized with `isSysFunc` and evaluated via `sysFuncEval`.
- External effects (creating geometry, animation) are done through the `pExternAPI` interface — the interpreter itself does not render, it calls out to the provided API.

Debugging and iteration limits
- If a `pDebug` object is provided, rule invocations are counted in `pDebug.RulesCallsCounterMap`.
- Global and per-rule iteration limits (`iters_max`) are enforced to avoid infinite recursion/loops; when reached, the rule exits and control returns to the caller, with stacks adjusted accordingly.

Behavioral notes / important points
- Expressions are a list of elements where the first element often identifies the construct (e.g., `if`, `set`, a variable name, a rule name).
- Sub-expressions return either a value or nil; nil means the sub-expression performed stateful actions only (no value to substitute).
- `return` statements inside subtrees allow returning a computed value up the call chain.
- System rules are side-effectful and call the extern API; user rules are more about state evolution and can call other rules.

If you want, I can:
- Walk through a short example AST and show step-by-step how `executeTree` evaluates it.
- Show the exact lines implementing a particular behavior (e.g., rule iteration limits, property incrementing), with direct code pointers.
