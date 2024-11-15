Sample Led Board with 8 LEDs
============================

Variant 1
```logi
board LedBoard {
    components [
        Led led1,
        Led led2,
        Led led3,
        Led led4,
        Led led5,
        Led led6,
        Led led7,
        Led led8
        
        Resistor 330 R1,
        Resistor 330 R2,
        Resistor 330 R3,
        Resistor 330 R4,
        Resistor 330 R5,
        Resistor 330 R6,
        Resistor 330 R7,
        Resistor 330 R8
        
        Wire w1,
        Wire w2,
        Wire w3,
        Wire w4,
        Wire w5,
        Wire w6,
        Wire w7,
        Wire w8
    ]
    
    connections [
        led1.a -> R1.a,
        led1.k -> R1.k,
        R1.a -> w1.a,
        R1.k -> w1.k,
        
        led2.a -> R2.a,
        led2.k -> R2.k,
        R2.a -> w2.a,
        R2.k -> w2.k,
        
        led3.a -> R3.a,
        led3.k -> R3.k,
        R3.a -> w3.a,
        R3.k -> w3.k,
        
        led4.a -> R4.a,
        led4.k -> R4.k,
        R4.a -> w4.a,
        R4.k -> w4.k,
        
        led5.a -> R5.a,
        led5.k -> R5.k,
        R5.a -> w5.a,
        R5.k -> w5.k,
        
        led6.a -> R6.a,
        led6.k -> R6.k,
        R6.a -> w6.a,
        R6.k -> w6.k,
        
        led7.a -> R7.a,
        led7.k -> R7.k,
        R7.a -> w7.a,
        R7.k -> w7.k,
        
        led8.a -> R8.a,
        led8.k -> R8.k,
        R8.a -> w8.a,
        R8.k -> w8.k
    ]
    
    power [
        w1.a,
        w2.a,
        w3.a,
        w4.a,
        w5.a,
        w6.a,
        w7.a,
        w8.a
    ]
}
```