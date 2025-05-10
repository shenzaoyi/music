package my_utils

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

func ExtractLastNumber(rawurl string) (int, error) {
	parsedUrl, err := url.Parse(rawurl)
	if err != nil {
		return 0, err
	}
	path := parsedUrl.Path
	path = strings.Trim(path, "/") // 去掉前后斜杠
	// 如果路径中有多级，用最后一段
	parts := strings.Split(path, "/")
	lastPart := parts[len(parts)-1]
	num, err := strconv.Atoi(lastPart)
	if err != nil {
		return 0, fmt.Errorf("最后部分不是数字: %v", err)
	}
	return num, nil
}
