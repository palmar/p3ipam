package main

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"strings"

	"github.com/fatih/color"
	bolt "go.etcd.io/bbolt"
)

var dbPath = "p3ipam.db"

func showHelp() {
	fmt.Println("p3ipam - commands: add, del, list, ping")
}

func openDB() (*bolt.DB, error) {
	return bolt.Open(dbPath, 0600, nil)
}

func handleAdd(args []string) {
	if len(args) == 0 {
		fmt.Println("Usage: p3ipam add subnet <cidr> name <name> | ip <ip> subnet <name>")
		return
	}
	if args[0] == "subnet" && len(args) == 4 && args[2] == "name" {
		cidr := args[1]
		name := args[3]
		db, err := openDB()
		if err != nil {
			fmt.Println("DB error:", err)
			return
		}
		defer db.Close()
		db.Update(func(tx *bolt.Tx) error {
			b, _ := tx.CreateBucketIfNotExists([]byte("subnets"))
			return b.Put([]byte(name), []byte(cidr))
		})
		color.Green("Added subnet %s => %s", name, cidr)
	} else if args[0] == "ip" && len(args) == 4 && args[2] == "subnet" {
		ip := args[1]
		subnet := args[3]
		db, err := openDB()
		if err != nil {
			fmt.Println("DB error:", err)
			return
		}
		defer db.Close()
		err = db.Update(func(tx *bolt.Tx) error {
			b, _ := tx.CreateBucketIfNotExists([]byte("ips"))
			return b.Put([]byte(ip), []byte(subnet))
		})
		if err == nil {
			color.Green("Added ip %s to subnet %s", ip, subnet)
		}
	} else {
		fmt.Println("Invalid add command")
	}
}

func handleDel(args []string) {
	if len(args) == 0 {
		fmt.Println("Usage: p3ipam del subnet <name> | ip <ip>")
		return
	}
	db, err := openDB()
	if err != nil {
		fmt.Println("DB error:", err)
		return
	}
	defer db.Close()
	if args[0] == "subnet" && len(args) == 2 {
		name := args[1]
		db.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("subnets"))
			if b != nil {
				b.Delete([]byte(name))
			}
			return nil
		})
		color.Red("Deleted subnet %s", name)
	} else if args[0] == "ip" && len(args) == 2 {
		ip := args[1]
		db.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("ips"))
			if b != nil {
				b.Delete([]byte(ip))
			}
			return nil
		})
		color.Red("Deleted ip %s", ip)
	} else {
		fmt.Println("Invalid del command")
	}
}

func handleList(args []string) {
	db, err := openDB()
	if err != nil {
		fmt.Println("DB error:", err)
		return
	}
	defer db.Close()
	if len(args) == 0 {
		db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("subnets"))
			if b == nil {
				fmt.Println("No subnets")
				return nil
			}
			b.ForEach(func(k, v []byte) error {
				fmt.Printf("%s => %s\n", k, v)
				return nil
			})
			return nil
		})
	} else if args[0] == "subnet" && len(args) == 2 {
		name := args[1]
		db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("subnets"))
			if b != nil {
				cidr := b.Get([]byte(name))
				if cidr != nil {
					fmt.Printf("%s => %s\n", name, cidr)
				}
			}
			return nil
		})
	} else {
		fmt.Println("Invalid list command")
	}
}

func pingHost(ip string) bool {
	cmd := exec.Command("ping", "-c", "1", "-w", "1", ip)
	err := cmd.Run()
	return err == nil
}

func ipsInCIDR(cidr string) []string {
	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil
	}
	var ips []string
	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
		ips = append(ips, ip.String())
	}
	if len(ips) > 0 {
		ips = ips[1 : len(ips)-1]
	}
	return ips
}

func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

func handlePing(args []string) {
	if len(args) == 0 {
		fmt.Println("Usage: p3ipam ping <ip|subnet|name>")
		return
	}
	target := args[0]
	var hosts []string
	db, _ := openDB()
	if db != nil {
		db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("subnets"))
			if b != nil {
				cidr := b.Get([]byte(target))
				if cidr != nil {
					hosts = ipsInCIDR(string(cidr))
				}
			}
			return nil
		})
		db.Close()
	}
	if len(hosts) == 0 {
		if strings.Contains(target, "/") {
			hosts = ipsInCIDR(target)
		} else {
			hosts = []string{target}
		}
	}
	for _, h := range hosts {
		if pingHost(h) {
			color.Green("%s alive", h)
		} else {
			color.Red("%s unreachable", h)
		}
	}
}

func main() {
	if len(os.Args) < 2 {
		showHelp()
		return
	}
	cmd := os.Args[1]
	switch cmd {
	case "add":
		handleAdd(os.Args[2:])
	case "del":
		handleDel(os.Args[2:])
	case "list":
		handleList(os.Args[2:])
	case "ping":
		handlePing(os.Args[2:])
	default:
		showHelp()
	}
}
