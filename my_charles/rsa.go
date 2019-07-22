package my_charles

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"io/ioutil"
	"math/big"
	rd "math/rand"
	"os"
	"time"
)

var RootCa *x509.Certificate
var RootKey *rsa.PrivateKey
var RootCaContents []byte

func init() {
	rd.Seed(time.Now().UnixNano())
	var err error
	RootCa, RootKey, err = Parse("/home/peterq/dev/env/cert/ca.crt", "/home/peterq/dev/env/cert/ca.nopass.key")
	if err != nil {
		panic(err)
	}
	RootCaContents, err = ioutil.ReadFile("/home/peterq/dev/env/cert/ca.crt")
	if err != nil {
		panic(err)
	}
}

func writeToBuff(buffer *bytes.Buffer, Type string, p []byte) error {
	var b = &pem.Block{Bytes: p, Type: Type}
	return pem.Encode(buffer, b)
}

func genCertWithCa(info CertInformation) (certBuff, keyBuff *bytes.Buffer, err error) {
	certBuff = bytes.NewBuffer([]byte{})
	keyBuff = bytes.NewBuffer([]byte{})

	Crt := newCertificate(info)
	Key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return
	}

	var buf []byte
	if RootCa == nil || RootKey == nil {
		//创建自签名证书
		buf, err = x509.CreateCertificate(rand.Reader, Crt, Crt, &Key.PublicKey, Key)
	} else {
		//使用根证书签名
		buf, err = x509.CreateCertificate(rand.Reader, Crt, RootCa, &Key.PublicKey, RootKey)
	}
	if err != nil {
		return
	}

	err = writeToBuff(certBuff, "CERTIFICATE", buf)
	if err != nil {
		return
	}
	certBuff.Write(RootCaContents)

	buf = x509.MarshalPKCS1PrivateKey(Key)
	err = writeToBuff(keyBuff, "PRIVATE KEY", buf)
	return
}

type CertInformation struct {
	Country            []string
	Organization       []string
	OrganizationalUnit []string
	EmailAddress       []string
	Province           []string
	Locality           []string
	CommonName         string
	CrtName, KeyName   string
	IsCA               bool
	Names              []pkix.AttributeTypeAndValue
}

func CreateCRT(RootCa *x509.Certificate, RootKey *rsa.PrivateKey, info CertInformation) error {
	Crt := newCertificate(info)
	Key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}

	var buf []byte
	if RootCa == nil || RootKey == nil {
		//创建自签名证书
		buf, err = x509.CreateCertificate(rand.Reader, Crt, Crt, &Key.PublicKey, Key)
	} else {
		//使用根证书签名
		buf, err = x509.CreateCertificate(rand.Reader, Crt, RootCa, &Key.PublicKey, RootKey)
	}
	if err != nil {
		return err
	}

	err = write(info.CrtName, "CERTIFICATE", buf)
	if err != nil {
		return err
	}

	buf = x509.MarshalPKCS1PrivateKey(Key)
	return write(info.KeyName, "PRIVATE KEY", buf)
}

//编码写入文件
func write(filename, Type string, p []byte) error {
	File, err := os.Create(filename)
	defer File.Close()
	if err != nil {
		return err
	}
	var b = &pem.Block{Bytes: p, Type: Type}
	return pem.Encode(File, b)
}

func Parse(crtPath, keyPath string) (rootcertificate *x509.Certificate, rootPrivateKey *rsa.PrivateKey, err error) {
	rootcertificate, err = ParseCrt(crtPath)
	if err != nil {
		return
	}
	rootPrivateKey, err = ParseKey(keyPath)
	return
}

func ParseCrt(path string) (*x509.Certificate, error) {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	p := &pem.Block{}
	p, buf = pem.Decode(buf)
	return x509.ParseCertificate(p.Bytes)
}

func ParseKey(path string) (*rsa.PrivateKey, error) {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	p, buf := pem.Decode(buf)
	return x509.ParsePKCS1PrivateKey(p.Bytes)
}

func newCertificate(info CertInformation) *x509.Certificate {
	return &x509.Certificate{
		SerialNumber: big.NewInt(rd.Int63()),
		Subject: pkix.Name{
			Country:            info.Country,
			Organization:       info.Organization,
			OrganizationalUnit: info.OrganizationalUnit,
			Province:           info.Province,
			CommonName:         info.CommonName,
			Locality:           info.Locality,
			ExtraNames:         info.Names,
		},
		NotBefore:             time.Now(),                   //证书的开始时间
		NotAfter:              time.Now().AddDate(20, 0, 0), //证书的结束时间
		BasicConstraintsValid: true,                         //基本的有效性约束
		IsCA:                  info.IsCA,                    //是否是根证书
		ExtKeyUsage: []x509.ExtKeyUsage{
			x509.ExtKeyUsageServerAuth,
			x509.ExtKeyUsageClientAuth,
			x509.ExtKeyUsageCodeSigning,
			x509.ExtKeyUsageEmailProtection,
			x509.ExtKeyUsageIPSECEndSystem,
			x509.ExtKeyUsageIPSECTunnel,
			x509.ExtKeyUsageIPSECUser,
			x509.ExtKeyUsageTimeStamping,
			x509.ExtKeyUsageOCSPSigning,
			x509.ExtKeyUsageMicrosoftServerGatedCrypto,
			x509.ExtKeyUsageNetscapeServerGatedCrypto,
			x509.ExtKeyUsageMicrosoftCommercialCodeSigning,
			x509.ExtKeyUsageMicrosoftKernelCodeSigning,
		}, //证书用途
		KeyUsage: x509.KeyUsageDigitalSignature |
			x509.KeyUsageContentCommitment |
			x509.KeyUsageKeyEncipherment |
			x509.KeyUsageDataEncipherment |
			x509.KeyUsageKeyAgreement |
			x509.KeyUsageCertSign |
			x509.KeyUsageCRLSign |
			x509.KeyUsageEncipherOnly |
			x509.KeyUsageDecipherOnly,
		//KeyUsage:       x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		EmailAddresses: info.EmailAddress,
		DNSNames:       []string{info.CommonName},
	}
}
