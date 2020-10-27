package dnslb

import (
	"fmt"
	"testing"
)

func TestDNS(t *testing.T) {
	dns := NewDNS(Config{})
	dns.Add("1", "12")
	dns.Add("2", "13")
	dns.Add("3", "14")
	dns.Add("4", "15")
	dns.Add("5", "16")

	fmt.Println("no deletion")
	for i := 0; i < 7; i++ {
		fmt.Println(dns.Next())
	}

	dns.Del("1")
	dns.Del("3")
	dns.Del("5")

	fmt.Println("with deletion")
	for i := 0; i < 8; i++ {
		fmt.Println(dns.Next())
	}

	dns.Add("8", "111")
	dns.Add("9", "234")
	dns.Add("10", "435")

	fmt.Println("with addition")
	for i := 0; i < 8; i++ {
		fmt.Println(dns.Next())
	}
}
