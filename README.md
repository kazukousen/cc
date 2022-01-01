
## Grammars

```
program      = funcDecl*
decl         = declspec declarator ("{" compoundStmt | ";")
funcDecl     = compoundStmt
declaration  = declspec declarator ("=" expr)? ("," declarator ("=" expr)?)*)? ";"
declspec     = "int" | "char"
declarator   = "*"* ident type-suffix
type-suffix  = "(" func-params | "[" num "]" type-suffix | ε
func-params  = (param ("," param)*)? ")"
param        = declspec declarator
stmt         = expr ";" | "{ compoundStmt | returnStmt | ifStmt | whileStmt | forStmt
compoundStmt = (declaration | stmt)* "}"
returnStmt   = "return" expr ";"
ifStmt       = "if" "(" expr ")" stmt ("else" stmt)?
whileStmt    = "while" "(" expr ")" stmt
forStmt      = "for" "(" expr? ";" expr? ";" expr? ")" stmt
expr         = assign
assign       = equality ("=" assign)?
equality     = relational ("==" relational | "!=" relational)*
relational   = add ("<" add | "<=" add | ">" add | ">=" add)*
add          = mul ("+" mul | "-" mul)*
mul          = unary ("*" unary | "/" unary)*
unary        = ("+" | "-" | "*" | "&") unary | postfix
postfix      = primary ("[" expr "]")*
primary      = "(" expr ")" | "sizeof" unary | ident func-args? | num
func-args    = "(" (assign ("," assign)*)? ")"
```
