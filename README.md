
## Grammars

```
program    = stmt*
stmt       = expr ";" | "{ stmt* "}" | returnStmt | ifStmt
returnStmt = "return" expr ";"
ifStmt     = "if" "(" expr ")" stmt ("else" stmt)?
whileStmt  = "while" "(" expr ")" stmt
expr       = assign
assign     = equality ("=" assign)?
equality   = relational ("==" relational | "!=" relational)*
relational = add ("<" add | "<=" add | ">" add | ">=" add)*
add        = mul ("+" mul | "-" mul)*
mul        = unary ("*" unary | "/" unary)*
unary      = ("+"|"-")? unary | primary
primary    = num | ident | "(" expr ")"
```
