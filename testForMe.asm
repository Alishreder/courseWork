.model small
.DATA
    HEX db 12h
a db 127
a1 db 100001b
a2 db 0fah
dec1 dw 15
b1 dw 111111111111000b
b2 dw 0c123h
    bin dd 1100101b
c dd 35
c1 dd 1111111111111110000b
c2 dd 12345678h

.CODE
 Je label1
label1:
 Je label1
    Std
    Div dword ptr [edx + esi]
Div byte ptr dS:[edx + esi]
Div byte ptr sS:[edx + esi]
Div dword ptr [esp + esi]
Div dword ptr sS:[edx + esi]
Div dword ptr eS:[edx + esi]
    jmp label2
    je label2
    Imul ax
Imul al
Imul eax
    Add ebx, ebx
    And bx, [ebp + esi]
    Cmp dec1, 13h
 jmp label2
label2:
 jmp label2
    Shl ebx, 10b
    Mov byte ptr[ecx + eax], 0bh
Mov word ptr es:[ecx + eax], 0bh
Mov word ptr[ecx + eax], 0bbbh
    Je label1
    jmp label1
div a
div cs:a
div c
END