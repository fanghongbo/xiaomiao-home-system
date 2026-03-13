package password

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

// Generate 创建一个满足条件的随机密码
func Generate(length int) (string, error) {
	var (
		password []byte
		allChars string
	)

	if length < 4 {
		return "", fmt.Errorf("password length must be at least 4")
	}

	// 确保至少包含一个小写字母、一个大写字母、一个数字和一个特殊字符
	password = append(password, getRandomChar(lowerLetters))
	password = append(password, getRandomChar(upperLetters))
	password = append(password, getRandomChar(digits))
	password = append(password, getRandomChar(specialChars))

	// 填充剩余部分
	allChars = lowerLetters + upperLetters + digits + specialChars
	for i := 4; i < length; i++ {
		password = append(password, getRandomChar(allChars))
	}

	// 手动实现Fisher-Yates洗牌算法来打乱顺序
	shuffle(password)

	return string(password), nil
}

// getRandomChar 使用 crypto/rand 获取随机字符
func getRandomChar(charSet string) byte {
	var (
		m *big.Int
		n *big.Int
	)

	m = big.NewInt(int64(len(charSet)))
	n, _ = rand.Int(rand.Reader, m)

	return charSet[n.Int64()]
}

// shuffle 使用Fisher-Yates洗牌算法打乱切片元素的顺序
func shuffle(data []byte) {
	for i := len(data) - 1; i > 0; i-- {
		j, _ := rand.Int(rand.Reader, big.NewInt(int64(i+1)))
		data[i], data[j.Int64()] = data[j.Int64()], data[i]
	}
}
