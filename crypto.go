package main

import "crypto/rand"
import "crypto/block"
import "crypto/aes"
import "bytes"
import "crypto/rsa"
import "encoding/base64"
import "crypto/sha256"
import "fmt"
import "big"
//import "strings"


func GenKey() (key *rsa.PrivateKey){
  key,_ = rsa.GenerateKey(rand.Reader, 2048)
  
  return key
}

func Sign(message []byte) []byte {
  
  key := GetKey()
  hasher := sha256.New()
  hasher.Write(message)
  hash := hasher.Sum()
  sig,err := rsa.SignPKCS1v15(rand.Reader, &key, rsa.HashSHA256, hash)
  if err != nil {
    fmt.Println(err)
  }
  return sig
}

func Verify(message []byte, sig []byte, name string) bool{
  fmt.Println(name)
  mod := Base64Decode([]byte(name))
  modbig := big.NewInt(0)
  modbig.SetBytes(mod) 
  key := &rsa.PublicKey{E:3,N:modbig}
  hasher := sha256.New()
  hasher.Write(message)
  hash := hasher.Sum()

  test := rsa.VerifyPKCS1v15(key, rsa.HashSHA256, hash, sig)
  if test != nil {
    return false
  }
  return true
}

func Encrypt(data []byte, key *rsa.PublicKey)(iv []byte, etext []byte, ckey []byte){ 
  var IVBuf [32]byte // create random 32 byte buffer for IV
  var EncodedText []byte
  rand.Read(IVBuf[0:])
  var SessionKey[32]byte
  rand.Read(IVBuf[0:]) // random 32byte session key
  var CryptedText []byte
  WriteBuf := new(bytes.Buffer)
  cipher,_ := aes.NewCipher(SessionKey[0:])
  BlockWriter := block.NewCBCEncrypter(cipher , IVBuf[0:], WriteBuf)
  BlockWriter.Write(data)
  EncryptedKey,_ := rsa.EncryptPKCS1v15(rand.Reader, key, SessionKey[0:])
  BlockWriter.Write(data) // Pass data through Encryptor
  WriteBuf.Read(CryptedText) // Read Encrypted data from buffer to byte slice
  //n := base64.URLEncoding.EncodedLen(len(CryptedText)))
  base64.URLEncoding.Encode(EncodedText, CryptedText)

  return IVBuf[0:], EncodedText, EncryptedKey
}

func Base64Encode(data []byte) []byte{
  var EncodedText [100000]byte
  n := base64.URLEncoding.EncodedLen(len(data))
  base64.URLEncoding.Encode(EncodedText[0:n], data)
  return EncodedText[0:n]
}

func Base64Decode(data []byte) []byte{
  var DecodedText [100000]byte
  n := base64.URLEncoding.DecodedLen(len(data))
  base64.URLEncoding.Decode(DecodedText[0:n], data)
  var i int; for i = len(DecodedText); i > 0 && DecodedText[i-1] == 0; i-- {}; 
  return DecodedText[:i]
}

func b64test(){
  key := GetKey()
  fmt.Println(key.PublicKey.N)
  
  modbytes := key.PublicKey.N.Bytes() // we want this byte value to re-create bigint
  //bytes := []byte{0x9, 0x14, 0x20}
  fmt.Println("bytes:",modbytes)
  
  
  encoded := Base64Encode(modbytes) //base64 encode modbytes, pass value == encoded to decode
  fmt.Println("encoded bytes:",string(encoded))
  fmt.Println("username:",GetUsername())
  original := Base64Decode(encoded) 
  fmt.Println("decoded bytes:",original)
  
  bigint := big.NewInt(0)
  bigint.SetBytes(original)
  fmt.Println(bigint)
  
  
}
