package secret

import (
	"github.com/infraboard/mcube/v2/crypto/cbc"
	"github.com/infraboard/mcube/v2/ioc/config/application"
)

// 加密
// 非对称加密: 加密密钥 (密钥对)
// 对称加密: 数据 (key: "xxx")
// 消息摘要(Hash)
// key 通过配置保存
// passwd   cbc.v1.abc.<cipherText>
func (s *Secret) Encrypt() error {
	cipherText, err := cbc.EncryptToString(s.Spec.Value, []byte(application.Get().EncryptKey))
	if err != nil {
		return err
	}
	s.Spec.Value = cipherText
	return nil
}

func (s *Secret) Decrypt() error {
	planText, err := cbc.DecryptFromString(s.Spec.Value, []byte(application.Get().EncryptKey))
	if err != nil {
		return err
	}
	s.Spec.Value = planText
	return nil
}

func (s *Secret) Desense() {
	// api secret
	// api****ret
	s.Spec.Value = s.Spec.Value[:3] + "****"
}
