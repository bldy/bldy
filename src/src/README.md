# src

`src` is a domain specific language designed for source manipulations.

## Variables
Variables can be declared by using the `let` keyword.
```
let x int
let y int = 123
let z = 1234
let s = "string"
let ns string
let na string = "asdf"
let Δ float64 = 1.0
let ε = Δ
```
## Types
Types can be declared by using the `type` keyword.

```
type MACH file
```
Union types can be declared like so.
```
type File {
    X
    Y
}
```
## Functions 
Functions can be defined by using the `func` keyword. 
Functions in `src` are very similar to functions in `Go`.
```
func C() {}
func compare(s string) {}
func A(s string) (int) {}
a = func() () {}
func y() () {
    func() () {

    }()
}
func t() () {
     a =func(s string) (){

     }
     a()
}
```
