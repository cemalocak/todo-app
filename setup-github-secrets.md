# 🔐 GitHub Secrets Kurulum Kılavuzu

## Repository Settings → Secrets and Variables → Actions

Aşağıdaki secrets'ları **"New repository secret"** ile ekleyin:

### 🌐 AWS CONFIGURATION
```
Secret Name: AWS_ACCESS_KEY_ID
Value: AKIA... (IAM'den aldığınız Access Key ID)
```

```
Secret Name: AWS_SECRET_ACCESS_KEY  
Value: ... (IAM'den aldığınız Secret Access Key)
```

```
Secret Name: AWS_REGION
Value: us-east-1 (EC2'nizi oluşturduğunuz region)
```

### 🖥️ EC2 CONFIGURATION
```
Secret Name: EC2_HOST
Value: 3.84.123.45 (EC2'nizin Public IP'si)
```

```
Secret Name: EC2_SSH_KEY
Value: 
-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEA...
(Tüm .pem dosyası içeriği)
...
-----END RSA PRIVATE KEY-----
```

### 🔗 ENVIRONMENT URLS
```
Secret Name: TEST_API_URL
Value: http://3.84.123.45:8080 (EC2_HOST + :8080)
```

```
Secret Name: PROD_API_URL  
Value: http://3.84.123.45 (EC2_HOST)
```

## ⚠️ ÖNEMLİ NOTLAR:

1. **EC2_SSH_KEY**: Tam .pem dosyası içeriğini kopyalayın (header/footer dahil)
2. **IP Adresleri**: http:// prefix'i ile yazın
3. **Region**: EC2 instance'ınızın bulunduğu AWS region
4. **Boşluk karakterleri**: Value'larda baş/sondaki boşlukları kaldırın

## ✅ KONTROL LİSTESİ:

- [ ] AWS_ACCESS_KEY_ID
- [ ] AWS_SECRET_ACCESS_KEY
- [ ] AWS_REGION
- [ ] EC2_HOST
- [ ] EC2_SSH_KEY (tam .pem içeriği)
- [ ] TEST_API_URL
- [ ] PROD_API_URL

Bu secrets'lar tanımlandıktan sonra GitHub Actions otomatik deployment yapabilir! 