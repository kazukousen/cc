
## Grammars

```
program    = funcDecl*
funcDecl   = declspec declarator stmt
declspec   = "int"
declarator = "*"* ident ("(" funcParams? ")")?
stmt       = expr ";" | "{ stmt* "}" | returnStmt | ifStmt | whileStmt | forStmt
returnStmt = "return" expr ";"
ifStmt     = "if" "(" expr ")" stmt ("else" stmt)?
whileStmt  = "while" "(" expr ")" stmt
forStmt    = "for" "(" expr? ";" expr? ";" expr? ")" stmt
expr       = assign
assign     = equality ("=" assign)?
equality   = relational ("==" relational | "!=" relational)*
relational = add ("<" add | "<=" add | ">" add | ">=" add)*
add        = mul ("+" mul | "-" mul)*
mul        = unary ("*" unary | "/" unary)*
unary      = ("+" | "-" | "*" | "&") unary | primary
primary    = num | ident ("(" funcParams? ")")? | "(" expr ")"
funcParams = assign ("," assign)*
```
