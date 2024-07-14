package tls

import ("fmt")

func hexdump(buf []byte) {
	for i, b := range buf[:] {
		if (i%16 == 0) {
			// print last 16 bytes as characters
			if i != 0 {
				fmt.Printf("  ")
				for j := i-16; j < i; j++ {
					if buf[j] >= 32 && buf[j] <= 126 {
						fmt.Printf("%c", buf[j])
					} else {
						fmt.Printf(".")
					}
				}
			}

			fmt.Printf("\n%04X  ", i)
		}
		fmt.Printf("%02X ", b)
	}
}
