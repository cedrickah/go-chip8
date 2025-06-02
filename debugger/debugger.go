package debugger

import (
	"fmt"
	"io"
)

func print_binary(w io.Writer, num uint8) {
	for c := 8; c > 0; c-- {
		if ((c)%4==0) {
            fmt.Fprintf(w, " ");
        }
        fmt.Fprintf(w, "%d", (num>>c)&1);	
	}
}

func Print_mem(w io.Writer, mem []uint8, from uint16, to uint16) {
	for i := from; i < to; i++ {
		fmt.Fprintf(w, "mem[%d|0x%04X]=", i, i);
        print_binary(w, mem[i]);
        fmt.Fprintf(w, "\n");
	}
}


func Print_instr(w io.Writer, instr uint8) {
		print_binary(w, instr);
        fmt.Fprintf(w, "\n");
}
