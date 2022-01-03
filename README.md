
## Grammars

```
program      = funcDecl*
decl         = declspec declarator ("{" funcDecl | varDecl)
varDecl      = ("," declarator)* ";"
funcDecl     = compoundStmt
declaration  = declspec declarator ("=" expr)? ("," declarator ("=" expr)?)*)? ";"
declspec     = "int" | "char" | struct-decl
declarator   = "*"* ident type-suffix
struct-decl  = "{" (declspec declarator ("," declarator)* ";")* "}"
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
postfix      = primary ("[" expr "]" | "." ident | "->" ident)*
primary      = "(" expr ")" | "sizeof" unary | ident func-args? | num | str
func-args    = "(" (assign ("," assign)*)? ")"
```
