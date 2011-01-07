package main

import "crypto/rand"
import "crypto/block"
import "crypto/aes"
import "bytes"
import "encoding/base64"

import "crypto/rsa"

func GenKey() (key *rsa.PrivateKey){
  key,_ = rsa.GenerateKey(rand.Reader, 2048)
  return key
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
  base64.URLEncoding.Encode(EncodedText, CryptedText)

  return IVBuf[0:], EncodedText, EncryptedKey
}
