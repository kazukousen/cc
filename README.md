
## Grammars

```
program      = funcDecl*
funcDecl     = declspec declarator "{" compoundStmt
declaration  = declspec declarator ("=" expr)? ("," declarator ("=" expr)?)*)? ";"
declspec     = "int"
declarator   = "*"* ident ("(" (declspec declarator ("," declspec declarator)*)? ")")?
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
unary        = ("+" | "-" | "*" | "&") unary | primary
primary      = num | ident ("(" (assign ("," assign)*)? ")")? | "(" expr ")"
```
