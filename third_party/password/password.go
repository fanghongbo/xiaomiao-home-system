package password

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"hash"
	"math/rand"
	"time"
	"unsafe"
)

func pbkdf2(password, salt []byte, iter, keyLen int, h func() hash.Hash) []byte {
	var (
		prf       hash.Hash
		hashLen   int
		numBlocks int
		buf       [4]byte
		dk        []byte
		U         []byte
	)

	prf = hmac.New(h, password)
	hashLen = prf.Size()
	numBlocks = (keyLen + hashLen - 1) / hashLen
	dk = make([]byte, 0, numBlocks*hashLen)
	U = make([]byte, hashLen)

	for block := 1; block <= numBlocks; block++ {
		// N.B.: || means concatenation, ^ means XOR
		// for each block T_i = U_1 ^ U_2 ^ ... ^ U_iter
		// U_1 = PRF(password, salt || uint(i))
		prf.Reset()
		prf.Write(salt)
		buf[0] = byte(block >> 24)
		buf[1] = byte(block >> 16)
		buf[2] = byte(block >> 8)
		buf[3] = byte(block)
		prf.Write(buf[:4])
		dk = prf.Sum(dk)
		T := dk[len(dk)-hashLen:]
		copy(U, T)

		// U_n = PRF(password, U_(n-1))
		for n := 2; n <= iter; n++ {
			prf.Reset()
			prf.Write(U)
			U = U[:0]
			U = prf.Sum(U)
			for x := range U {
				T[x] ^= U[x]
			}
		}
	}
	return dk[:keyLen]
}

func build(base string, length int, factor int, b []byte) []byte {
	var r *rand.Rand
	r = rand.New(rand.NewSource(time.Now().Unix()))
	for i, cacheValue, remain := length-factor, r.Int63(), letterIdxMax; i >= 0; {
		var idx int
		if remain == 0 {
			cacheValue, remain = r.Int63(), letterIdxMax
		}
		if idx = int(cacheValue & letterIdxMask); idx < len(base) {
			b[i] = base[idx]
			i--
		}
		cacheValue >>= letterIdxBits
		remain--
	}
	return b
}

func NewSalt(length int) (string, error) {
	var (
		r      *rand.Rand
		b      []byte
		factor int
	)
	// assign elements from 4 kinds of base elements
	r = rand.New(rand.NewSource(time.Now().Unix()))
	factor = length / 4
	b = make([]byte, length)
	b = build(lowerLetters, length, 1, b)
	b = build(upperLetters, length, factor*2, b)
	b = build(specialChars, length, factor*3, b)
	b = build(digits, length, factor*4, b)
	// shuffle
	rand.Shuffle(len(b), func(i, j int) {
		i = r.Intn(length)
		b[i], b[j] = b[j], b[i]
	})
	return *(*string)(unsafe.Pointer(&b)), nil
}

func New(password string, salt string) (string, error) {
	return hex.EncodeToString(pbkdf2([]byte(password), []byte(salt), 10000, 50, sha256.New)), nil
}
