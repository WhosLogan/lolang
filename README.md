# Lolang
A compiler and virtual machine I built because I was bored.

Code is first written in human-readable `.lo` files and compiled by the compiler into `.loc` files which are then read 
and executed by the virtual machine.

Note that this is just a personal project and probably has quite a few unchecked bugs.

## Project Structure
- `/compiler` - The compiler of Lolang, responsible for turning .lo files into compiled .loc files by generating an
instruction set that can be read by the virtual machine.
- `/disasm` - A disassembler I made to help with debugging, takes .loc files and prints their functions and 
instructions.
- `/runtime` - The entire runtime for the language. Has standard library features as well as the entire virtual machine.
- `/shared` - A shared module that contains structures and basic code used by all of the other projects in the workspace.

## How to use
Start by compiling your `.lo` file:
```bash
compiler HelloWorld.lo HelloWorld.loc
```

Then by running it:
```bash
runtime HelloWorld.loc
```

## Example

Below is an example of Lolang as well as the compiled output:
```lolang
lo printTimes(string msg, int n) {
    for (int i = 0; i < n; i = i + 1) {
        println(msg);
    }
    return;
}

lo main() {
    printTimes("looping...", 3);
    return;
}

```

## Custom Types (Structs)

Lolang supports user-defined struct types. Here's an example:

```lolang
struct Point {
    int x;
    int y;
}

struct Rectangle {
    Point topLeft;
    Point bottomRight;
}

lo main() {
    Point p = new Point { x: 10, y: 20 };
    printInt(p.x);
    printInt(p.y);
    
    p.x = 30;
    printInt(p.x);
    
    return;
}
```

### Struct Features
- **Definition**: Use `struct Name { ... }` to define a new type with fields
- **Instantiation**: Use `new StructName { field: value, ... }` to create instances
- **Field Access**: Use dot notation `instance.field` to read fields
- **Field Assignment**: Use `instance.field = value` to update fields
- **Nested Structs**: Structs can contain other structs as fields (e.g., `rect.topLeft.x`)
- **Type Checking**: Full type checking for struct fields and assignments

### New Struct Opcodes
- `NewObj` - Creates a new struct instance of the specified type
- `LdFld` - Loads a field value from a struct onto the stack
- `StFld` - Stores a value from the stack into a struct field

## Compiled Output Example

The first example (printTimes) compiles to:
```
Fn #1000000: main
Instr #0: Opcode: LdStr Operand: looping...
Instr #1: Opcode: Ldc8 Operand: 3
Instr #2: Opcode: Call Operand: 1000001
Instr #3: Opcode: Ret Operand: nil

Fn #1000001: printTimes
Instr #0: Opcode: StLoc Operand: 0
Instr #1: Opcode: StLoc Operand: 1
Instr #2: Opcode: Ldc8 Operand: 0
Instr #3: Opcode: StLoc Operand: 2
Instr #4: Opcode: LdLoc Operand: 2
Instr #5: Opcode: LdLoc Operand: 1
Instr #6: Opcode: Clt Operand: nil
Instr #7: Opcode: Bf Operand: 15
Instr #8: Opcode: LdLoc Operand: 0
Instr #9: Opcode: Call Operand: 1
Instr #10: Opcode: LdLoc Operand: 2
Instr #11: Opcode: Ldc8 Operand: 1
Instr #12: Opcode: Add Operand: nil
Instr #13: Opcode: StLoc Operand: 2
Instr #14: Opcode: Br Operand: 4
Instr #15: Opcode: Ret Operand: nil
```

## Notes
- The standard library, while incredibly small, is insanely easy to add to as the virtual machine can just call Go 
functions if they're from the standard library.
- Not every instruction is currently used due to the limited code generator. However the virtual machine is fairly 
well-built out.

## Todo
I'm sure I'll eventually get bored and come back to work on this project, so here's a quick list of things I 
should improve on:

- [x] Add custom types that can be defined by the user
- [ ] More primitive types
- [ ] More robust testing suite
- [ ] Some sort of error handling (probably done similarly to go with tuples)
- [ ] Rework the way arguments are passed so they are not just pushed onto the stack of the target function
- [ ] Adding on to above, rework the entire calling convention to be more future-proof and support modules
- [ ] Base value type is still too limiting
- [ ] Type and TypeCode needs to be reworked, it was semi half-assed

## Credits
- Microsoft. Cannot stress this enough, DotNet's MSIL was a **huge** help when I initially learned how to write virtual 
machines (especially stack based ones) back when I was younger. Without Microsoft this project would cease to exist.
- My professor who let me do this project instead of an essay. She's amazing and without her I would've never done this.