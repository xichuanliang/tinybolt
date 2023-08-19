package main

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"
	"unsafe"
)

type person struct {
	age    int
	weight int
	height int
}

// 加载配置文件
func loadConfig(filePath string) (map[string]string, error) {
	config := make(map[string]string)
	data, err := ioutil.ReadFile(filePath) // 读取配置文件内容
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(data), "\n") // 按行拆分配置项
	for _, line := range lines {
		parts := strings.SplitN(line, "=", 2) // 将配置项拆分为键值对
		if len(parts) == 2 {
			config[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1]) // 添加配置项到map
		}
	}
	return config, nil
}

// 替换变量
func replaceVariables(yamlContent string, config map[string]string) string {
	re := regexp.MustCompile(`\${([^}]*)}`)
	replacedContent := re.ReplaceAllStringFunc(yamlContent, func(match string) string {
		varName := match[2 : len(match)-1]
		value, exists := config[varName]
		if exists {
			return value
		}
		return match
	})
	return replacedContent
}

func main() {
	yamlFilePath := "clusterctl-template.yaml" // YAML 文件路径
	configFilePath := "input.conf"             // 配置文件路径
	outputFilePath := "output.yaml"            // 输出文件路径

	config, err := loadConfig(configFilePath) // 加载配置文件
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		return
	}

	yamlContent, err := ioutil.ReadFile(yamlFilePath) // 读取 YAML 文件内容
	if err != nil {
		fmt.Printf("Error reading YAML file: %v\n", err)
		return
	}

	updatedYAMLContent := replaceVariables(string(yamlContent), config) // 替换变量

	err = ioutil.WriteFile(outputFilePath, []byte(updatedYAMLContent), 0644) // 将更新的内容写入输出文件
	if err != nil {
		fmt.Printf("Error writing output YAML file: %v\n", err)
		return
	}

	fmt.Println("Variable replacement completed. Updated YAML content saved to output.yaml")

	pos := unsafe.Offsetof(person{}.height)
	fmt.Println(pos)
	p := &person{
		age:    1,
		weight: 12,
		height: 14,
	}
	ppos := unsafe.Offsetof(p.height)

	fmt.Println(ppos)
	fmt.Println("-------------------------------")
	a := "aa"
	dst := (*[10]byte)(unsafe.Pointer(&a))[:5]
	fmt.Println(dst)
	fmt.Println("we have to save our life")
}
