package main

//func main() {
//	in, err := os.Open("./nginx.conf")
//	if err != nil {
//		fmt.Println("open file fail:", err)
//		os.Exit(-1)
//	}
//	defer in.Close()
//	out, err := os.OpenFile("nginx.conf", os.O_RDWR|os.O_CREATE, 0766)
//	if err != nil {
//		fmt.Println("Open write file fail:", err)
//		os.Exit(-1)
//	}
//	defer out.Close()
//	br := bufio.NewReader(in)
//	index := 0
//
//	str := "server {\n    listen  80;\n    server_name     b.com;\n    location / {\n        proxy_pass      http://localhost:10608;\n    }\n}\n#<?>\n"
//	file_lines := make([]string, 0)
//	posi := 0
//	for {
//		line, _, err := br.ReadLine()
//		if err == io.EOF {
//			break
//		}
//		file_lines = append(file_lines, string(line))
//		if strings.Contains(string(line), "#<?>") {
//			posi = index
//		}
//		index++
//	}
//	file_lines[posi] = strings.ReplaceAll(file_lines[posi], "#<?>", str)
//	for _, line := range file_lines {
//		fmt.Println(line)
//		_, err = out.WriteString(line + "\n")
//		if err != nil {
//			fmt.Println("write to file fail:", err)
//			os.Exit(-1)
//		}
//	}
//
//}
//
//func main() {
//	in, err := os.Open("./nginx.conf")
//	if err != nil {
//		fmt.Println("open file fail:", err)
//		os.Exit(-1)
//	}
//	defer in.Close()
//	out, err := os.OpenFile("nginx.conf", os.O_RDWR|os.O_CREATE, 0766)
//	if err != nil {
//		fmt.Println("Open write file fail:", err)
//		os.Exit(-1)
//	}
//	defer out.Close()
//	br := bufio.NewReader(in)
//	index := 0
//
//	file_lines := make([]string, 0)
//	posi := 0
//	for {
//		line, _, err := br.ReadLine()
//		if err == io.EOF {
//			break
//		}
//		file_lines = append(file_lines, string(line))
//		if strings.Contains(string(line), "http://localhost:10608") {
//			posi = index
//		}
//		index++
//	}
//
//	file_lines = append(file_lines[:posi-4], file_lines[posi+3:]...)
//
//	//file_lines[posi] = strings.ReplaceAll(file_lines[posi], "#<?>", str)
//	for _, line := range file_lines {
//		fmt.Println(line)
//		//_, err = out.WriteString(line + "\n")
//		//if err != nil {
//		//	fmt.Println("write to file fail:", err)
//		//	os.Exit(-1)
//		//}
//	}
//
//}
