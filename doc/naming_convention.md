**Naming convention**  
There are multiple languages used - Go/Python/Typescript/Rust. The goal is to maintain a simple universal naming scheme across all languages. For some languages this scheme is not standard, but having it be consistent across all of the code (including the shared symbol names) has its benefits in readibility and correctness.  
Rules:
- snake_case
- function argument names beging with "p_" to easily indicate right away where the value is coming from (outside the function, or from internal scope).
- if values/variables are of generic type, such as string/float/int/list/map/tuple, then their names should end with a postfix with a shorthand. If its a custom/user_defined type 
  then there is no posftix. this practice increases readibility and acts as local documentation, for which types are involved in a particular expression, either in dynamic languages,
  or in languages with type inferencers.  
  Type suffixes:
    - float  - "_f"  
    - int    - "_int"  
    - string - "_str"  
    - list   - "_lst"  
    - map    - "_map"  

**Code Style**
A single style is maintained across languages used in the implementation (**Go**,**Python**,**Typescript**) - even though the languages are different enough from each other. 
The focus is on basic functional language principles (of pure functions, high-level functions, closures). Functions should receive all the state that they operate on via their arguments (other then functions that work with external state - DB or external queries). Object orientation (objects holding state and methods operating on that state internally) is avoided as much as possible (even though it is the default idiomatic style of Go and Python). State/variable mutation still exists in various places, but the aim is to keep it to a minimum (constant runtime values would be a welcome feature in Go and Python). 


    
