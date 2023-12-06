DEF CurrentGameNum EQU $C000
DEF RunningCount EQU $C002
DEF NumRed EQU $C004
DEF NumGreen EQU $C006
DEF NumBlue EQU $C008
DEF GameNotPossible EQU $C00A

SECTION "Header", ROM0[$100]

Main:
  ld bc, InputData
  xor a
  ld d, a

LoadGameNum:
  ld a, [bc]
  inc bc
  cp a, $20 ; space
  jp nz, LoadGameNum ; keep looping til we hit a space
  ld a, [bc] ; write char after space to serial
  ld [CurrentGameNum], a
  inc bc
  ld a, [bc]
  inc bc
  cp a, $3A ; colon
  jp z, LoadColor; if char after first digit is a colon, move on to color
  inc bc ; TODO: store second game digit somewhere
  jp LoadColor

LoadColor:
  ld a, [bc]
  inc bc
  cp a, $3B ; semicolon
  jp z, EvalGame; TODO: implement this
  cp a, $A ; newline
  jp z, EndOfGame
  cp a, $39 ; check if its a digit
  jp nc, LoadColor
  cp a, $30
  jp c, LoadColor ; past here we know its a digit
  ld l, a
  inc bc
  ld a, [bc]
  cp a, $72 ; r for red
  jp z, LoadRed
  cp a, $62 ; b for blue
  jp z, LoadBlue
  cp a, $67 ; g for green
  jp z, LoadGreen
  jp LoadColor

LoadRed:
  ld a, l
  ld [NumRed], a
  jp LoadColor

LoadBlue:
  ld a, l
  ld [NumBlue], a
  jp LoadColor

LoadGreen:
  ld a, l
  ld [NumGreen], a
  jp LoadColor

; TODO: if game is possible, keep going to end of line. if not, new line?
; need to store state for line somehow?
EvalGame:
  call Debug
  ld a, [NumRed]
  cp a, $D ; 13 red TODO: doesn't work b/c we have 2-digit numbers :(
  jp c, LoadColor ; means 13 is greater than NumRed, so game possible
  ld a, [NumGreen]
  cp a, $E ; 14 green
  jp nc, LoadColor ; means 14 is greater than NumGreen, so game possible
  ld a, [NumBlue]
  cp a, $F ; 15 blue
  jp nc, LoadColor ; means 15 is greater than NumBlue, so game possible
  ld a, $1
  ld [GameNotPossible], a
  jp LoadColor

EndOfGame:
  jp Wait

Debug:
  ld a, $47 ; G
  ld [$ff01], a
  ld [$ff02], a
  ld a, $47 ; G
  ld [$ff01], a
  ld [$ff02], a
  ld a, [CurrentGameNum]
  ld [$ff01], a
  ld [$ff02], a
  ld a, $A ; newline
  ld [$ff01], a
  ld [$ff02], a

  ld a, $43 ; C
  ld [$ff01], a
  ld [$ff02], a
  ld a, [RunningCount]
  ld [$ff01], a
  ld [$ff02], a
  ld a, $A ; newline
  ld [$ff01], a
  ld [$ff02], a

  ld a, $52 ; R
  ld [$ff01], a
  ld [$ff02], a
  ld a, [NumRed]
  ld [$ff01], a
  ld [$ff02], a
  ld a, $A; newline
  ld [$ff01], a
  ld [$ff02], a

  ld a, $42; B
  ld [$ff01], a
  ld [$ff02], a
  ld a, [NumBlue]
  ld [$ff01], a
  ld [$ff02], a
  ld a, $A; newline
  ld [$ff01], a
  ld [$ff02], a

  ld a, $47; G
  ld [$ff01], a
  ld [$ff02], a
  ld a, [NumGreen]
  ld [$ff01], a
  ld [$ff02], a
  ld a, $A; newline
  ld [$ff01], a
  ld [$ff02], a
  ret

Wait:
  jp Wait

InputData: INCBIN "test1.data"
