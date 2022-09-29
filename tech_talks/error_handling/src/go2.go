func main() {
	handle err {
		log.Fatal(err)
	}

	hex := check ioutil.ReadAll(os.Stdin)
	data := check parseHexdump(string(hex))
	os.Stdout.Write(data)
}