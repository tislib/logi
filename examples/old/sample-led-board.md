Sample Led Board with 8 LEDs
============================

Variant 1
```logi
board SampleLedBoard {
    inputs [vcc, gnd, data, clock, latch]
    
    components [
        led1 Led(color=red)
        led2 Led(color=green)
        led3 Led(color=blue)
        led4 Led(color=yellow)
        led5 Led(color=white)
        led6 Led(color=orange)
        led7 Led(color=purple)
        led8 Led(color=black)
        
        shiftRegister ShiftRegister
        
        resistor1 Resistor(value=220)
        resistor2 Resistor(value=220)
        resistor3 Resistor(value=220)
        resistor4 Resistor(value=220)
        resistor5 Resistor(value=220)
        resistor6 Resistor(value=220)
        resistor7 Resistor(value=220)
        resistor8 Resistor(value=220)
    ]
    
    connections [
        from gnd to resistor1
        from gnd to resistor2
        from gnd to resistor3
        from gnd to resistor4
        from gnd to resistor5
        from gnd to resistor6
        from gnd to resistor7
        from gnd to resistor8
        
        from resistor1 to led1
        from resistor2 to led2
        from resistor3 to led3
        from resistor4 to led4
        from resistor5 to led5
        from resistor6 to led6
        from resistor7 to led7
        from resistor8 to led8
        
        from vcc to shiftRegister
        from data to shiftRegister
        from clock to shiftRegister
        from latch to shiftRegister
        
        from shiftRegister to led1
        from shiftRegister to led2
        from shiftRegister to led3
        from shiftRegister to led4
        from shiftRegister to led5
        from shiftRegister to led6
        from shiftRegister to led7
        from shiftRegister to led8
        
        from shiftRegister to resistor1
        from shiftRegister to resistor2
        from shiftRegister to resistor3
        from shiftRegister to resistor4
        from shiftRegister to resistor5
        from shiftRegister to resistor6
        from shiftRegister to resistor7
        from shiftRegister to resistor8
    ]

}
```