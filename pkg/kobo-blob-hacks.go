package pkg

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
)

func hacksType44(data []byte) (map[string]interface{}, error) {
	if bytes.Contains(data, []byte(pocketMime)) {
		extraDataReadingSecsZeros, _ := hex.DecodeString(extraDataReadingSecondsWithZeros)

		if bytes.Contains(data, extraDataReadingSecsZeros) {
			idx := bytes.Index(data, extraDataReadingSecsZeros)
			sp := idx + len(extraDataReadingSecsZeros) + 3 // hardcoded, I guess this is 1 padding byte and 2 bytes for length or type
			ep := sp + 4                                   // hardcoded 4 bytes for int32
			r := bytes.NewReader(data[sp:ep])

			var v int32
			if err := binary.Read(r, binary.LittleEndian, &v); err != nil { // TODO! hardcoded byte order !!!
				panic(err)
			}

			/*fmt.Println(string(data))
			fmt.Println(fmt.Sprintf("v = %d", v/100000))
			panic("got here")*/

			return map[string]interface{}{
				"ContentType":           pocketMime,
				extraDataReadingSeconds: v / 100000, // hardcoded 100000 division to seconds
			}, nil
		}
	}

	return map[string]interface{}{}, nil
}
